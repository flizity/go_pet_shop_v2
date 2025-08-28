package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
}

func StatusHandler(status, message string, code int, data any) Response {
	return Response{
		Status:  status,
		Message: message,
		Code:    code,
		Data:    data,
	}
}

func HealthCheck(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.status.HealthCheck"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Health check requested", slog.String("url", r.URL.String()))

		response := StatusHandler("success", "Service is healthy", http.StatusOK, map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"service":   "go_pet_shop",
			"version":   "1.0.0",
		})

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response)
	}
}

func ReadinessCheck(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.status.ReadinessCheck"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Readiness check requested", slog.String("url", r.URL.String()))

		response := StatusHandler("success", "Service is ready", http.StatusOK, map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"database":  "connected",
			"service":   "go_pet_shop",
		})

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response)
	}
}

func StatusInfo(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.status.StatusInfo"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Status info requested", slog.String("url", r.URL.String()))

		response := StatusHandler("success", "Service status information", http.StatusOK, map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"uptime":    time.Since(time.Now().Add(-time.Hour)).String(),
			"service":   "go_pet_shop",
			"version":   "1.0.0",
			"env":       "development",
			"endpoints": []string{
				"/health",
				"/ready",
				"/status",
				"/api/users",
				"/api/products",
				"/api/orders",
			},
		})

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, response)
	}
}
