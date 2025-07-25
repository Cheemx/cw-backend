package shared

import (
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func RequireAuth(authHeader string) (bool, error) {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return false, errors.New("authorization header missing")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return false, errors.New(err.Error())
	}
	return true, nil
}
