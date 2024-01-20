package server_utils

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
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

	if len(md["user_id"]) == 0 {
		return ctx, status.Error(codes.Unauthenticated, "jwt not found")
	}

	userId := md["user_id"][0]
	//token, err := utils.ValidateAccessToken(tokenString)
	//
	//if err != nil {
	//	if len(md["refresh_token"]) == 0 {
	//		return ctx, status.Error(codes.Unauthenticated, "invalid token, refresh token not found")
	//	}
	//
	//	refreshToken := md["refresh_token"][0]
	//	newAccessToken, err := utils.RefreshAccessToken(refreshToken)
	//	if err != nil {
	//		return ctx, status.Error(codes.Unauthenticated, "failed to refresh token")
	//	}
	//
	//	newContext := context.WithValue(ctx, "jwt", newAccessToken)
	//	return newContext, nil
	//}
	//
	//claims, ok := token.Claims.(jwt.MapClaims)
	//if !ok {
	//	return ctx, status.Error(codes.Unknown, "failed to extract claims")
	//}
	//
	//userID, ok := claims["uid"].(string)
	//if !ok {
	//	return ctx, status.Error(codes.Unauthenticated, "user ID not found in token")
	//}

	newContext := context.WithValue(ctx, "user_id", userId)
	return newContext, nil
}
