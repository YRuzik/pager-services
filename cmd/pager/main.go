package main

import (
	"crypto/tls"
	_ "embed"
	"flag"
	mux2 "github.com/gorilla/mux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net"
	"net/http"
	"pager-services/pkg/mongo_ops"
	"pager-services/pkg/server_utils"
	handlers "pager-services/pkg/sockets"
	"pager-services/pkg/transfers"
)

//go:embed certs/server.crt
var certTLS []byte

//go:embed certs/server.key
var keyTLS []byte

// fsdf
func init() {
	mongo_ops.InitMongoDB()
	transfers.PagerLockerInitialize()

	log.SetFlags(log.Lshortfile | log.LstdFlags)
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

	grpcAddress := "0.0.0.0:0"
	httpAddress := "0.0.0.0:4001"
	authAddress := "0.0.0.0:0"
	httpAuthAddress := "0.0.0.0:5001"

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

	tcpHttpAuthListener, listenerError := net.Listen("tcp", httpAuthAddress)
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

	hub := handlers.NewHub()

	go func() { server_utils.StartGrpcServer(tlsGrpcListener) }()
	go func() { server_utils.StartAuthServer(tlsAuthListener) }()
	go hub.Run()

	c := server_utils.Cors()

	proxy := server_utils.ProxyBuilder(tlsGrpcListener.Addr().String(), tlsConfig)
	authProxy := server_utils.ProxyBuilder(tlsAuthListener.Addr().String(), tlsConfig)

	mux := mux2.NewRouter()
	httpAuthMux := http.NewServeMux()

	server_utils.HandleHttpRoutes(mux, hub)

	http2Server := &http2.Server{}
	http1Server := &http.Server{Handler: h2c.NewHandler(c.Handler(server_utils.CreateGrpcWithHttpHandler(mux, proxy, false)), http2Server)}

	httpAuth2Server := &http2.Server{}
	http1AuthServer := &http.Server{Handler: h2c.NewHandler(c.Handler(server_utils.CreateGrpcWithHttpHandler(httpAuthMux, authProxy, true)), httpAuth2Server)}

	go func() {
		log.Print("[HTTPS AUTH SERVER] server listening on address: ", tcpHttpAuthListener.Addr().String())
		if httpServerError := http1AuthServer.Serve(tcpHttpAuthListener); httpServerError != nil {
			return
		}
	}()

	log.Print("[HTTPS SERVER] server listening on address: ", tlsHttpListener.Addr().String())
	if httpServerError := http1Server.Serve(tcpHttpListener); httpServerError != nil {
		return
	}

}
