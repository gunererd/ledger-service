package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap the response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// Call the next handler
		next.ServeHTTP(wrapped, r)
		
		// Log the request with structured data
		duration := time.Since(start)
		logger.Info("HTTP request",
			"method", r.Method,
			"path", r.RequestURI,
			"remote_addr", r.RemoteAddr,
			"status_code", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}