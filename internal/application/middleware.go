package application

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status     int
	error      string
	expression string
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		log.Printf("Входящий %s запрос на %s", r.Method, r.URL.Path)

		next.ServeHTTP(wrapped, r)

		logMessage := "Завершен %s %s - статус: %d, длительность: %v"
		logArgs := []interface{}{
			r.Method,
			r.URL.Path,
			wrapped.status,
			time.Since(start),
		}

		if wrapped.status != http.StatusOK {
			logMessage += ", ошибка: %s"
			if wrapped.expression != "" {
				logMessage += ", выражение: %s"
				logArgs = append(logArgs, wrapped.error, wrapped.expression)
			} else {
				logArgs = append(logArgs, wrapped.error)
			}
		}

		log.Printf(logMessage, logArgs...)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
