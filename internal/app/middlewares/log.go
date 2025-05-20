package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/aube/url-shortener/internal/logger"
)

// responseData holds information about the HTTP response.
type responseData struct {
	status int // HTTP status code
	size   int // Response size in bytes
}

// loggingResponseWriter wraps http.ResponseWriter to capture response details.
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write captures the response size while writing to the underlying ResponseWriter.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader captures the status code while writing to the underlying ResponseWriter.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// String provides a string representation of the response data.
func (r *loggingResponseWriter) String() {
	var buf bytes.Buffer
	buf.WriteString("Response:")
	buf.WriteString("Headers:")
	for k, v := range r.ResponseWriter.Header() {
		buf.WriteString(fmt.Sprintf("%s: %v", k, v))
	}
}

// LoggingMiddleware logs details about HTTP requests and responses.
// It captures:
// - Request method and URI
// - Response status code
// - Response size
// - Duration of the request
// - Content encoding and type headers
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		log.Info(
			"LoggingMiddleware",
			"status", responseData.status,
			"method", r.Method,
			"URI", r.RequestURI,
			"duration", duration,
			"size", responseData.size,
			"ce", r.Header.Get("Content-Encoding"),
			"ct", r.Header.Get("Content-Type"),
		)
	})
}
