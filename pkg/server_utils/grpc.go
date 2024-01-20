package server_utils

import (
	"context"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	pager_auth "pager-services/pkg/api/pager_api/auth"
	pager_chat "pager-services/pkg/api/pager_api/chat"
	pager_transfers "pager-services/pkg/api/pager_api/transfers"
	"pager-services/pkg/auth"
	"pager-services/pkg/chat_actions"
	"pager-services/pkg/transfers"
	"strings"
)

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

func StartGrpcServer(lis net.Listener) {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(AuthInterceptor), grpc.StreamInterceptor(AuthStreamInterceptor))
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

func StartAuthServer(lis net.Listener) {
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

func CreateGrpcWithHttpHandler(httpHand http.Handler, proxy httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if multiplexer.IsGrpcWebRequest(r) {
			multiplexer.ServeHTTP(w, r)
			return
		}
		if r.Method == "POST" && strings.HasPrefix(r.Header.Get("content-type"), "application/grpc") {
			proxy.ServeHTTP(w, r)
			return
		}
		httpHand.ServeHTTP(w, r)
	})
}
