package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Estrutura para armazenar limitadores por IP
var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
	rateLimit = rate.Limit(5)  // 5 requisições por segundo
	burstSize = 10             // Máximo de 10 requisições de uma vez
)

// Obtém ou cria um rate limiter para um IP específico
func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rateLimit, burstSize)
		visitors[ip] = limiter

		// Remover visitantes antigos para evitar consumo excessivo de memória
		go func(ip string) {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			delete(visitors, ip)
			mu.Unlock()
		}(ip)
	}
	return limiter
}

// Middleware para Rate Limiting
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getVisitor(ip)

		// Verifica se pode processar a requisição
		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
