package httpapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		requestID := GetRequestID(r.Context())
		duration := time.Since(start)

		log.Printf("[%s] %s %s %d %s",
			requestID,
			r.Method,
			r.URL.Path,
			rec.statusCode,
			duration,
		)
	})
}

func GetRequestID(ctx context.Context) string {
	v := ctx.Value(requestIDKey)
	if id, ok := v.(string); ok {
		return id
	}
	return "unknown"
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
