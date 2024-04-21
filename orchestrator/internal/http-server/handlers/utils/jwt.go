package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func ParseJWTToken(r *http.Request) (string, string, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// В этой функции необходимо вернуть секретный ключ, используемый для подписи токена
		return []byte("your-256-bit-secret"), nil
	})
	if err != nil {
		return "", "", err
	}
	fmt.Println("333333333")
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", ErrNoJWT
	}
	fmt.Println("4444444444444444")
	if token.Valid {
		fmt.Println("4444444444444444")
		userName, _ok := claims["name"].(string)
		if !_ok {
			return "", "", ErrNoJWT
		}
		userId, _ok := claims["id"].(string)
		if !_ok {
			return "", "", ErrNoJWT
		}
		return userId, userName, nil
	} else {
		fmt.Println("555555555555555555")
		return "", "", err
	}
}
