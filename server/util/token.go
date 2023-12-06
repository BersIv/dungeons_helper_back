package util

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

type MyJWTClaims struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	IdAvatar int64  `json:"idAvatar"`
	jwt.RegisteredClaims
}

func GetIdFromHeader(r *http.Request) (int64, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("Authorization header is missing")
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return 0, errors.New("Invalid token format")
	}

	tokenString := splitToken[1]

	secret := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(tokenString, &MyJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*MyJWTClaims)
	if !ok || !token.Valid {
		return 0, errors.New("Invalid token")
	}

	return claims.Id, nil
}
