package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	common "pager-services/pkg/api/pager_api/common"
	"time"
)

var secretKey = []byte("testSecret")

func NewToken(uid primitive.ObjectID, identity string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = uid
	claims["identity"] = identity
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", MentorError("failed create new token", codes.Internal, &common.PagerError{
			Code:    common.PagerError_INTERNAL,
			Details: err.Error(),
		})
	}
	return tokenString, nil
}

func NewRefreshToken(uid primitive.ObjectID, identity string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = uid
	claims["identity"] = identity
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	//return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, status.Error(codes.Unknown, "unexpected signing method")
	//	}
	//
	//	return secretKey, nil
	//})
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unknown, "unexpected signing method:")
		}
		return secretKey, nil
	})
	switch {
	case token.Valid:
		return token, nil
	case errors.Is(err, jwt.ErrTokenExpired):
		return nil, status.Error(codes.Unauthenticated, "token expired")
	default:
		return nil, status.Error(codes.Unknown, "validate access token failed")
	}
}

func ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, MentorError("unexpected signing method", codes.Unknown, &common.PagerError{
				Code: common.PagerError_UNKNOWN,
			})
		}

		return secretKey, nil
	})
}

func RefreshAccessToken(refreshToken *jwt.Token) (string, error) {

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", MentorError(" failed to extract claims from refresh token", codes.Unauthenticated, &common.PagerError{
			Code: common.PagerError_UNAUTHENTICATED,
		})
	}

	uid, ok := claims["user_id"].(string)
	if !ok {
		return "", MentorError("user ID not found in refresh token", codes.Unauthenticated, &common.PagerError{
			Code: common.PagerError_UNAUTHENTICATED,
		})
	}
	identity, ok := claims["identity"].(string)
	objectUserId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return "", MentorError("failed get objectUserId", codes.Internal, &common.PagerError{
			Code:    common.PagerError_INTERNAL,
			Details: err.Error(),
		})
	}
	if !ok {
		return "", MentorError("user ID not found in refresh token", codes.Unauthenticated, &common.PagerError{
			Code: common.PagerError_UNAUTHENTICATED,
		})
	}
	newAccessToken, err := NewToken(objectUserId, identity, 5*time.Second)
	if err != nil {
		return "", MentorError("failed to generate new access token", codes.Internal, &common.PagerError{
			Code:    common.PagerError_INTERNAL,
			Details: err.Error(),
		})
	}

	return newAccessToken, nil
}
