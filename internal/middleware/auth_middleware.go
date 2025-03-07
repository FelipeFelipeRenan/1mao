package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// AuthMiddleware agora aceita roles e retorna um mux.MiddlewareFunc
func AuthMiddleware(roles []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Token não fornecido", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			secret := os.Getenv("JWT_SECRET")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Erro ao processar token", http.StatusUnauthorized)
				return
			}

			userID := uint(claims["user_id"].(float64))
			userRole := claims["role"].(string) // Supondo que a role esteja armazenada no token

			// Verificar se a role do usuário é permitida para acessar a rota
			allowed := false
			for _, role := range roles {
				if role == userRole {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Acesso negado", http.StatusForbidden)
				return
			}

			// Adiciona o userID ao contexto da requisição
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
