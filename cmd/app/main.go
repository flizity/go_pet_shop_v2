package main

import (
	"context"
	"go_pet_shop/internal/config"
	handlers "go_pet_shop/internal/delivery/http"
	"go_pet_shop/internal/lib/logger"
	kafkaRepo "go_pet_shop/internal/repository/kafka"
	"go_pet_shop/internal/repository/postgres"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)
	log.Info("starting go_pet_shop server",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.Address))

	repository, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer repository.Close()

	log.Info("database connection established")

	// KafkaProducer
	producer, err := kafkaRepo.NewKafkaProducer([]string{"localhost:9092"})
	if err != nil {
		log.Error("failed to connect to kafka", slog.Any("error", err))
		os.Exit(1)
	}
	defer producer.Close()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	router.Post("/api/send", handlers.SendToKafka(log, producer))

	// Status
	router.Get("/health", handlers.HealthCheck(log))

	// Order history
	router.Route("/api/history", func(r chi.Router) {
		r.Get("/{email}", handlers.GetUserOrderHistory(log, repository))
		r.Get("/popular", handlers.GetPopularProducts(log, repository))
	})

	// Orders
	router.Route("/api/orders", func(r chi.Router) {
		r.Post("/", handlers.CreateOrder(log, repository))
	})

	// Products
	router.Route("/api/products", func(r chi.Router) {
		r.Post("/", handlers.CreateProduct(log, repository))
	})

	// Transactions
	router.Route("/api/transactions", func(r chi.Router) {
		r.Post("/", handlers.PlaceOrder(log, repository))
	})

	// Users
	router.Route("/api/users", func(r chi.Router) {
		r.Post("/", handlers.CreateUser(log, repository))
	})

	router.Get("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "API is working", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server starting", slog.String("address", cfg.Address))

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed to start", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	log.Info("server started successfully")

	<-done
	log.Info("server stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("server stopped gracefully")
}
