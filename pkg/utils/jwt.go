package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return "", err
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
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unknown, "unexpected signing method")
		}

		return secretKey, nil
	})
}

func ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Error(codes.Unknown, "unexpected signing method")
		}

		return secretKey, nil
	})
}

func RefreshAccessToken(refreshTokenString string) (string, error) {
	refreshToken, err := ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "failed to extract claims from refresh token")
	}

	uid, ok := claims["user_id"].(string)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "user ID not found in refresh token")
	}
	identity, ok := claims["identity"].(string)
	objectUserId, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", status.Error(codes.Unauthenticated, "user ID not found in refresh token")
	}
	newAccessToken, err := NewToken(objectUserId, identity, 5*time.Minute)
	if err != nil {
		return "", status.Error(codes.Internal, "failed to generate new access token")
	}

	return newAccessToken, nil
}
