package server_utils

import (
	"crypto/tls"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	handlers "pager-services/pkg/sockets"
)

func HandleHttpRoutes(mux *http.ServeMux, hub *handlers.Hub) {
	mux.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		_, err := io.WriteString(responseWriter, "hello inreko practice")
		if err != nil {
			log.Print("error while write string ", err)
			return
		}
	})

	mux.HandleFunc("/ws/{userID}", func(responseWriter http.ResponseWriter, request *http.Request) {
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
