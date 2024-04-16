package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func ParseJWTToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// В этой функции необходимо вернуть секретный ключ, используемый для подписи токена
		return []byte("your-256-bit-secret"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// В данном блоке можно обрабатывать найденные утверждения из токена
		userId := claims["name"].(string)
		return userId, nil
	} else {
		return "", err
	}
}
