package server_utils

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"pager-services/pkg/utils"
)

func AuthStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if newContext, err := getNewContext(ss.Context()); err != nil {
		return err
	} else {
		return handler(srv, &serverStream{ss, newContext})
	}
}

func AuthInterceptor(ctx context.Context,
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

	if len(md["jwt"]) == 0 {
		return ctx, status.Error(codes.Unauthenticated, "invalid token, refresh token not found")
	}

	tokenString := md["jwt"][0]
	token, err := utils.ValidateAccessToken(tokenString)

	if err != nil {
		return ctx, status.Error(codes.Unauthenticated, "invalid token, refresh token not found")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx, status.Error(codes.Unknown, "failed to extract claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "user ID not found in token")
	}

	newContext := context.WithValue(ctx, "user_id", userID)
	return newContext, nil
}
