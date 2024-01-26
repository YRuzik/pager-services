package server_utils

import (
	"crypto/tls"
	mux2 "github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	handlers "pager-services/pkg/sockets"
)

func HandleHttpRoutes(mux *mux2.Router, hub *handlers.Hub) {
	mux.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		_, err := io.WriteString(responseWriter, "hello inreko practice")
		if err != nil {
			log.Print("error while write string ", err)
			return
		}
	})

	mux.HandleFunc("/ws/{userId}/", func(responseWriter http.ResponseWriter, request *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		userID := mux2.Vars(request)["userId"]

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		connection, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		log.Print("new connection: " + userID)

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
