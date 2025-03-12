package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const lokiEndpoint = "http://loki:3100/loki/api/v1/push"

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		logData := map[string]interface{}{
			"method":     r.Method,
			"path":       r.RequestURI,
			"ip":         r.RemoteAddr,
			"user_agent": r.UserAgent(),
			"duration":   time.Since(start).Milliseconds(), // Melhor para métricas
			"timestamp":  time.Now().Format(time.RFC3339Nano),
		}

		// Log formatado para o console
		log.WithFields(logData).Info("Acesso registrado")

		// Envio para Loki
		if err := sendLogToLoki(logData); err != nil {
			log.WithError(err).Error("Erro ao enviar log para Loki")
		}
	})
}

func sendLogToLoki(logData map[string]interface{}) error {
	logEntry, err := json.Marshal(logData)
	if err != nil {
		return fmt.Errorf("erro ao serializar logData: %w", err)
	}

	// Loki requer timestamps em nanosegundos como string
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	payload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": map[string]string{
					"app": "1mao",
				},
				"values": [][]string{
					{timestamp, string(logEntry)},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %w", err)
	}

	resp, err := http.Post(lokiEndpoint, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição ao Loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("loki retornou status inesperado: %d", resp.StatusCode)
	}

	return nil
}
