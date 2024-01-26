package server_utils

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"log"
	common "pager-services/pkg/api/pager_api/common"
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
		return ctx, utils.MentorError("md not found", codes.Unauthenticated, &common.PagerError{
			Code: common.PagerError_UNKNOWN,
		})
	}

	if len(md["jwt"]) == 0 {
		return ctx, utils.MentorError("invalid token, refresh token not found", codes.Unauthenticated, &common.PagerError{
			Code: common.PagerError_UNAUTHENTICATED,
		})
	}

	tokenString := md["jwt"][0]
	token, err := utils.ValidateAccessToken(tokenString)

	if err != nil {
		return ctx, utils.MentorError("invalid token, refresh token not found", codes.Unauthenticated, &common.PagerError{
			Code:    common.PagerError_UNAUTHENTICATED,
			Details: err.Error(),
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx, utils.MentorError("failed to extract claims", codes.Internal, &common.PagerError{
			Code:    common.PagerError_INTERNAL,
			Details: err.Error(),
		})
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return ctx, utils.MentorError("user ID not found in token", codes.Unauthenticated, &common.PagerError{
			Code:    common.PagerError_UNAUTHENTICATED,
			Details: err.Error(),
		})
	}

	newContext := context.WithValue(ctx, "user_id", userID)
	return newContext, nil
}
