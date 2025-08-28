package logger

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func LoggingMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("HTTP request:",
			slog.String("method", r.Method),
			slog.String("url", r.URL.String()),
		)
		next.ServeHTTP(w, r)
	})
}

func CustomLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			log.Info("HTTP response",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Duration("duration", time.Since(start)),
				slog.String("remote", r.RemoteAddr),
			)
		})
	}
}
