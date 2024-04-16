package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func ParseJWTToken(r *http.Request) (string, string, error) {
	fmt.Println("11111111111111111111")
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer")
	fmt.Println("2222222222222222")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// В этой функции необходимо вернуть секретный ключ, используемый для подписи токена
		return []byte("your-256-bit-secret"), nil
	})
	fmt.Println("33333333333333333")
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userName := claims["name"].(string)
		userId := claims["id"].(string)
		fmt.Println("4444444444444444")
		return userId, userName, nil
	} else {
		fmt.Println("555555555555555555")
		return "", "", err
	}
}
