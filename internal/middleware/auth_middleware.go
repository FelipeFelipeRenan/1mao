package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == ""{
			http.Error( w,"Token não fornecido", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func (token *jwt.Token)(interface{}, error)  {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid{
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		claims, ok:= token.Claims.(jwt.MapClaims)
		if !ok{
			http.Error(w, "Erro ao processar token", http.StatusUnauthorized)
			return
		}

		userID := uint(claims["user_id"].(float64))

		ctx:= context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	
}