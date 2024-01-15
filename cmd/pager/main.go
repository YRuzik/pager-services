package main

import (
	"crypto/tls"
	_ "embed"
	"flag"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	pager_chat "pager-services/pkg/api/pager_api/chat"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/chat_actions"
	"pager-services/pkg/mongo_ops"
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

	tcpGrpcListener, listenerError := net.Listen("tcp", grpcAddress)
	if listenerError != nil {
		log.Fatalf("failed to listen: %v", listenerError)
	}

	tcpHttpListener, listenerError := net.Listen("tcp", httpAddress)
	if listenerError != nil {
		log.Fatalf("failed to listen: %v", listenerError)
	}

	tlsConfig, loadCredsError := loadTLSCredentials()
	if loadCredsError != nil {
		log.Fatalf("loadCreds error")
	}

	tlsHttpListener := tls.NewListener(tcpHttpListener, tlsConfig)
	tlsGrpcListener := tls.NewListener(tcpGrpcListener, tlsConfig)

	go func() { startGrpcServer(tlsGrpcListener) }()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowCredentials: true,
		Debug:            true,
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

	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(c.Handler(createGrpcWithHttpHandler(httpMux, proxy)), http2Server)}

	log.Print("[HTTPS SERVER] server listening on address: ", tlsHttpListener.Addr().String())
	if httpServerError := http1Server.Serve(tcpHttpListener); httpServerError != nil {
		return
	}
}

func startGrpcServer(lis net.Listener) {
	grpcServer := grpc.NewServer()
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
