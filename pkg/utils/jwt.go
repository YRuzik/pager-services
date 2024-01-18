package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(uid string, identity string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = uid
	claims["identity"] = identity
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte("testSecret"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
