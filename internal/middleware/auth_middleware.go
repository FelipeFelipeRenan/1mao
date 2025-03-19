package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const UserContextKey ContextKey = "user"

func AuthMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenHeader := r.Header.Get("Authorization")
			if tokenHeader == "" {
				http.Error(w, "Token não encontrado", http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(tokenHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				http.Error(w, "Token mal formatado", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenParts[1], func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			role := claims["role"].(string)
			log.Println("✅ Token válido para:", role)

			// Verifica se o papel do usuário está na lista permitida
			allowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Acesso negado", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
