package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"flag"
	"github.com/gorilla/websocket"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	pager_auth "pager-services/pkg/api/pager_api/auth"
	pager_chat "pager-services/pkg/api/pager_api/chat"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/auth"
	"pager-services/pkg/chat_actions"
	"pager-services/pkg/mongo_ops"
	handlers "pager-services/pkg/sockets"
	"pager-services/pkg/transfers"
	"strings"
)

//go:embed certs/server.crt
var certTLS []byte

//go:embed certs/server.key
var keyTLS []byte

func init() {
	mongo_ops.InitMongoDB()
}

type grpcMultiplexer struct {
	*grpcweb.WrappedGrpcServer
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}

var multiplexer grpcMultiplexer

func getRoot(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "hello inreko practice")
	if err != nil {
		log.Print("error while write string ", err)
		return
	}

}

func loadTLSCredentials() (*tls.Config, error) {
	certificate, err := tls.X509KeyPair(certTLS, keyTLS)

	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.NoClientCert,
		InsecureSkipVerify: true,
	}

	return config, nil
}

func main() {
	flag.Parse()
	//ctx := context.Background()

	grpcAddress := "localhost:0"
	httpAddress := "localhost:4001"
	authAddress := "localhost:5001"

	tcpGrpcListener, listenerError := net.Listen("tcp", grpcAddress)
	if listenerError != nil {
		log.Fatalf("failed to listen: %v", listenerError)
	}

	tcpHttpListener, listenerError := net.Listen("tcp", httpAddress)
	if listenerError != nil {
		log.Fatalf("failed to listen: %v", listenerError)
	}

	tcpAuthListener, listenerError := net.Listen("tcp", authAddress)
	if listenerError != nil {
		log.Fatalf("failed to listen: %v", listenerError)
	}

	tlsConfig, loadCredsError := loadTLSCredentials()
	if loadCredsError != nil {
		log.Fatalf("loadCreds error")
	}

	tlsHttpListener := tls.NewListener(tcpHttpListener, tlsConfig)
	tlsGrpcListener := tls.NewListener(tcpGrpcListener, tlsConfig)
	tlsAuthListener := tls.NewListener(tcpAuthListener, tlsConfig)

	go func() { startGrpcServer(tlsGrpcListener) }()
	go func() { startAuthServer(tlsAuthListener) }()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "https"
			req.URL.Host = tcpGrpcListener.Addr().String()
		},
		ErrorLog: log.New(log.Writer(), "", 0),
		Transport: &http.Transport{
			ForceAttemptHTTP2:  true,
			DisableCompression: true,
			TLSClientConfig:    tlsConfig,
		},
		FlushInterval: -1,
	}

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", getRoot)
	hub := handlers.NewHub()
	go hub.Run()
	httpMux.HandleFunc("/ws/{userID}", func(responseWriter http.ResponseWriter, request *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		// Reading username from request parameter
		userID := request.Header.Get("userId")

		// Upgrading the HTTP connection socket connection
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		connection, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		handlers.CreateNewSocketUser(hub, connection, userID)

	})

	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(c.Handler(createGrpcWithHttpHandler(httpMux, proxy)), http2Server)}

	log.Print("[HTTPS SERVER] server listening on address: ", tlsHttpListener.Addr().String())
	if httpServerError := http1Server.Serve(tcpHttpListener); httpServerError != nil {
		return
	}
}

func startGrpcServer(lis net.Listener) {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor), grpc.StreamInterceptor(authStreamInterceptor))
	reflection.Register(grpcServer)
	RegisterGrpcServices(grpcServer)
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	multiplexer = grpcMultiplexer{
		grpcWebServer,
	}
	log.Print("[GRPC SERVER] server listening on address: ", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		return
	}
}

func startAuthServer(lis net.Listener) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pager_auth.RegisterAuthServiceServer(grpcServer, auth.PagerAuth{})
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	multiplexer = grpcMultiplexer{
		grpcWebServer,
	}
	log.Print("[AUTH SERVER] server listening on address: ", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		return
	}
}

func RegisterGrpcServices(registrar grpc.ServiceRegistrar) {
	pager_chat.RegisterChatActionsServer(registrar, &chat_actions.PagerChat{})
	pager_transfers.RegisterPagerStreamsServer(registrar, &transfers.PagerStreams{})
}

func createGrpcWithHttpHandler(httpHand http.Handler, proxy httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if multiplexer.IsGrpcWebRequest(r) {
			multiplexer.ServeHTTP(w, r)
		}
		if r.Method == "POST" && strings.HasPrefix(r.Header.Get("content-type"), "application/grpc") {
			proxy.ServeHTTP(w, r)
			return
		}
		httpHand.ServeHTTP(w, r)
	})
}

func authStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if newContext, err := getNewContext(ss.Context()); err != nil {
		return err
	} else {
		return handler(srv, &serverStream{ss, newContext})
	}
}

// INTERCEPTOR
func authInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	log.Printf("\nRequest - Method: %s\t \nError: %v\n",
		info.FullMethod)
	if info.FullMethod == "/com.niokr.api.PollActions/RecalculateFitage" {
		handl, err := handler(ctx, req)
		log.Printf("\nRequest - Method: %s\t \nError: %v\n",
			info.FullMethod,
			err)

		return handl, err
	} else if newContext, err := getNewContext(ctx); err != nil {
		return nil, err
	} else {
		handl, err := handler(newContext, req)
		log.Printf("\nRequest - Method: %s\t \nError: %v\n",
			info.FullMethod,
			err)

		return handl, err
	}
}

func getNewContext(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {

		return ctx, status.Error(codes.Unauthenticated, "md not found")
	}
	if len(md["user_id"]) == 0 {
		//logging.Default.Error("jwt doesn't exists", zap.Any("error", "md: "+fmt.Sprintf("%+v", md)))
		return ctx, status.Error(codes.Unauthenticated, "jwt not found")
	}
	tokenString := md["user_id"][0]
	return context.WithValue(ctx, "user_id", tokenString), nil
	//if token, err := firebaseAuth.VerifyIDToken(ctx, tokenString); err != nil {
	//	return ctx, status.Error(codes.Unauthenticated, "invalid token")
	//} else {
	//	userId, ok := token.Claims["user_id"].(string)
	//	if !ok {
	//		return ctx, status.Error(codes.Unauthenticated, "id doesnt exists")
	//	} else {
	//		newContext := context.WithValue(ctx, "user_id", userId)
	//		return newContext, nil
	//	}
	//}
}
