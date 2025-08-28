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

type Products interface {
	CreateProduct(product models.Product) (int, error)
	GetProductByID(id int) (models.Product, error)
	GetAllProducts() ([]models.Product, error)
	UpdateProduct(product models.Product) error
	DeleteProduct(id int) error
}

func CreateProduct(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.CreateProduct"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request id", middleware.GetReqID(r.Context())),
		)

		log.Info("create new product", slog.String("url", r.URL.String()))

		var product models.Product
		if err := render.DecodeJSON(r.Body, &product); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := products.CreateProduct(product); err != nil {
			log.Error("failed to create product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Product created successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Product created successfully"})
	}
}

func GetProductByID(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.GetProductByID"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Error("failed to get product by id")
			http.Error(w, "Product ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("failed to parse id parameter", slog.Any("error", err))
			http.Error(w, "Product ID is invalid", http.StatusBadRequest)
			return
		}

		product, err := products.GetProductByID(id)
		if err != nil {
			log.Error("failed to get product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Info("Product retrieved successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, product)
	}
}

func GetAllProducts(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.GetAllProducts"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Getting all products", slog.String("url", r.URL.String()))

		allProducts, err := products.GetAllProducts()
		if err != nil {
			log.Error("failed to get all products", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("All products retrieved successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, allProducts)
	}
}

func UpdateProduct(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.UpdateProduct"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Updating product", slog.String("url", r.URL.String()))

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Error("empty id")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid id", slog.String("id", idStr), slog.Any("error", err))
			http.Error(w, "invalid product id", http.StatusBadRequest)
			return
		}

		var product models.Product
		if err := render.DecodeJSON(r.Body, &product); err != nil {
			log.Error("failed to decode request body", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		product.ID = id

		if err := products.UpdateProduct(product); err != nil {
			log.Error("failed to update product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Product updated successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Product updated successfully"})
	}
}

func DeleteProduct(log *slog.Logger, products Products) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.products.DeleteProduct"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("Deleting product", slog.String("url", r.URL.String()))

		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			log.Error("empty id")
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid id", slog.String("id", idStr), slog.Any("error", err))
			http.Error(w, "invalid product id", http.StatusBadRequest)
			return
		}

		if err := products.DeleteProduct(id); err != nil {
			log.Error("failed to delete product", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Product deleted successfully", slog.String("url", r.URL.String()))

		render.JSON(w, r, map[string]string{"status": "Product deleted successfully"})
	}
}
