package handlers

import (
	models "go_pet_shop/internal/domain"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Histories interface {
	GetUserOrderHistory(email string) ([]models.OrderDetail, error)
	GetPopularProducts() ([]models.PopularProduct, error)
}

func GetUserOrderHistory(log *slog.Logger, histories Histories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.histories.GetUserOrderHistory"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		email := chi.URLParam(r, "email")
		log.Info("Fetching order history for user", slog.String("email", email))

		history, err := histories.GetUserOrderHistory(email)
		if err != nil {
			log.Error("failed to get user order history", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, history)
	}
}
func GetPopularProducts(log *slog.Logger, histories Histories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.histories.GetPopularProducts"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Fetching popular products", slog.String("url", r.URL.String()))

		products, err := histories.GetPopularProducts()
		if err != nil {
			log.Error("failed to get popular products", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, products)
	}
}
