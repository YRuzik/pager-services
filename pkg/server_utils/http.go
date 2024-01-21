package server_utils

import (
	"crypto/tls"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	handlers "pager-services/pkg/sockets"
)

func HandleHttpRoutes(mux *http.ServeMux) {
	manager := handlers.NewManager()
	mux.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		_, err := io.WriteString(responseWriter, "hello inreko practice")
		if err != nil {
			log.Print("error while write string ", err)
			return
		}
	})

	mux.HandleFunc("/ws", manager.ServeWS)
}

func ProxyBuilder(grpcAddress string, tlsConfig *tls.Config) httputil.ReverseProxy {
	return httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "https"
			req.URL.Host = grpcAddress
		},
		ErrorLog: log.New(log.Writer(), "proxy error ", 0),
		Transport: &http.Transport{
			ForceAttemptHTTP2:  true,
			DisableCompression: true,
			TLSClientConfig:    tlsConfig,
		},
		FlushInterval: -1,
	}
}

func Cors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
}
