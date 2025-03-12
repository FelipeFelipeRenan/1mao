package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		// Log formatado com metodo, rota, tempo de resposta e ip
		log.Println("------------------------------------------------------------")
		log.Printf("[%s] %s | %s |%s |%v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			r.UserAgent(),
			time.Since(start),
		)
		log.Println("------------------------------------------------------------")
	})
}
