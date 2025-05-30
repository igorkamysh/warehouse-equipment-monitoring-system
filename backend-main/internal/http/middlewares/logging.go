package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Info("handle request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int64("duration_ns", int64(time.Since(start))),
		)
	})
}
