package middleware

import (
	"errors"
	"strings"

	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	if rw.status == 0 { // Evita múltiplas chamadas
		rw.status = code
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK // Se WriteHeader não foi chamado, assume-se 200
	}
	return rw.ResponseWriter.Write(b)
}

var cb *gobreaker.CircuitBreaker

func init() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "1Mao Circuit Breaker",
		MaxRequests: 5,                // Permite 5 tentativas antes de abrir
		Interval:    10 * time.Second, // Tempo de reset após falha
		Timeout:     30 * time.Second, // Tempo de eséra antes de abrir novamente
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 2 // Abre após 3 falhas consecutivas
		},
	})
}

func CircuitBreakerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ignorando o websocket
		if strings.HasPrefix(r.URL.Path, "/ws/") {
			next.ServeHTTP(w, r)
			return
		}
		rw := &responseWriterWrapper{ResponseWriter: w}

		_, err := cb.Execute(func() (interface{}, error) {
			next.ServeHTTP(rw, r)

			if rw.status >= 500 {
				log.Println("❌ Circuit Breaker registrou erro:", rw.status)
				return nil, errors.New("erro detectado na resposta")
			}

			return nil, nil
		})

		if err != nil {
			log.Println("🚨 Circuit Breaker BLOQUEOU a requisição!")
			http.Error(w, "Serviço temporariamente indisponível", http.StatusServiceUnavailable)
		}
	})
}
