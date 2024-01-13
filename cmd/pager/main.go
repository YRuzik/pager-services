package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"flag"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"pager-services/pkg/transfers"
	"pager-services/pkg/utils"
	"strings"
)

//go:embed certs/server.crt
var certTLS []byte

//go:embed certs/server.key
var keyTLS []byte

func init() {
	transfers.InitMongoDB()
}

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
	ctx := context.Background()

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
	http1Server := &http.Server{Handler: h2c.NewHandler(createGrpcWithHttpHandler(httpMux, proxy), http2Server)}

	log.Print("[HTTPS SERVER] server listening on address: ", tlsHttpListener.Addr().String())
	if httpServerError := http1Server.Serve(tlsHttpListener); httpServerError != nil {
		return
	}
	transfers.ReadStream(ctx, transfers.Client.Database("test_streams").Collection("transfers"), "test")
}

func startGrpcServer(lis net.Listener) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	utils.RegisterGrpcServices(grpcServer)
	log.Print("[GRPC SERVER] server listening on address: ", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		return
	}
}

func createGrpcWithHttpHandler(httpHand http.Handler, proxy httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && strings.HasPrefix(r.Header.Get("content-type"), "application/grpc") {
			proxy.ServeHTTP(w, r)
			return
		}
		httpHand.ServeHTTP(w, r)
	})
}
