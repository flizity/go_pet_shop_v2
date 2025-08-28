package handlers

import (
	models "go_pet_shop/internal/domain"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Transactions interface {
	PlaceOrder(userEmail string, items []models.OrderItem) (orderID int, err error)
}

func PlaceOrder(log *slog.Logger, transactions Transactions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.transactions.PlaceOrder"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Creating new order", slog.String("url", r.URL.String()))

		var req struct {
			Email string             `json:"email"`
			Items []models.OrderItem `json:"items"`
		}

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		}

		orderID, err := transactions.PlaceOrder(req.Email, req.Items)
		if err != nil {
			log.Error("failed to place order", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		log.Info("Order placed successfully", slog.Int("order_id", orderID))
		render.JSON(w, r, map[string]int{"order_id": orderID})
	}
}
