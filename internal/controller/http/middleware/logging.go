package middleware

import (
	"net/http"
	"time"

	"github.com/YelzhanWeb/uno-spicchio/pkg/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		logger.HTTP(
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}
