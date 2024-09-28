package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Withlogging logs every incoming request
func WithLogging(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := log.With(
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("addr", r.RemoteAddr),
		)

		start := time.Now()

		defer func() {
			entry.Info(
				"Requset",
				slog.Duration("duration", time.Since(start)),
			)
		}()

		next.ServeHTTP(w, r)
	})
}
