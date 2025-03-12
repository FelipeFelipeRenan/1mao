package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init(){
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		// Log formatado com metodo, rota, tempo de resposta e ip
		log.WithFields(logrus.Fields{
			"method": r.Method,
			"path": r.RequestURI,
			"ip": r.RemoteAddr,
			"user_agent": r.UserAgent(),
			"duration": time.Since(start).String(),
		}).Info("Acesso registrado")
	})
}
