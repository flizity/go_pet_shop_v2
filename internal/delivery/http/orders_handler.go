package handlers

import (
	models "go_pet_shop/internal/domain"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Orders interface {
	CreateOrder(order models.Order) (int, error)
	AddOrderItem(orderItem models.OrderItem) error
	GetOrderByID(id int) (models.Order, error)
	GetOrdersByUserEmail(email string) ([]models.Order, error)
	GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error)
}

func CreateOrder(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.CreateOrder"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Creating new order", slog.String("url", r.URL.String()))

		var order models.Order
		if err := render.DecodeJSON(r.Body, &order); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := orders.CreateOrder(order)
		if err != nil {
			log.Error("failed to create order", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Created new order", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]int{"id": id})

	}
}

func AddOrderItem(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.AddOrderItem"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Adding order item", slog.String("url", r.URL.String()))

		var orderItem models.OrderItem
		if err := render.DecodeJSON(r.Body, &orderItem); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := orders.AddOrderItem(orderItem); err != nil {
			log.Error("failed to add order item", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Order item added successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Order item added successfully"})
	}
}

func GetOrderByID(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.GetOrderByID"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Getting order by ID", slog.String("url", r.URL.String()))

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Error("missing id parameter")
			http.Error(w, "missing id parameter", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid id parameter", slog.String("id", idStr), slog.Any("error", err))
			http.Error(w, "invalid id parameter", http.StatusBadRequest)
			return
		}

		order, err := orders.GetOrderByID(id)
		if err != nil {
			log.Error("failed to get order", slog.Int("id", id), slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Info("Order retrieved successfully", slog.Int("id", id))

		render.JSON(w, r, order)
	}
}

func GetOrdersByUserEmail(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.GetOrdersByUserEmail"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		email := r.URL.Query().Get("email")

		if email == "" {
			log.Error("missing email parameter")
			http.Error(w, "missing email parameter", http.StatusBadRequest)
			return
		}
		log.Info("Getting orders by user email", slog.String("url", r.URL.String()))

		orders, err := orders.GetOrdersByUserEmail(email)
		if err != nil {
			log.Error("failed to get orders by user email", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Info("Order retrieved successfully", slog.String("url", r.URL.String()))
		render.JSON(w, r, orders)
	}
}

func GetOrderItemsByOrderID(log *slog.Logger, orders Orders) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.orders.GetOrderItemsByOrderID"

		log := log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		orderIDStr := chi.URLParam(r, "orderID")
		if orderIDStr == "" {
			log.Error("missing id parameter")
			http.Error(w, "missing id parameter", http.StatusBadRequest)
			return
		}
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			log.Error("invalid id parameter", slog.Any("error", err))
			http.Error(w, "invalid id parameter", http.StatusBadRequest)
			return
		}
		orderItems, err := orders.GetOrderItemsByOrderID(orderID)
		if err != nil {
			log.Error("failed to get order items", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Info("Order items retrieved successfully", slog.String("url", r.URL.String()))
		render.JSON(w, r, orderItems)
	}
}
