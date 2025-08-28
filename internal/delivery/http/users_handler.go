package handlers

import (
	"encoding/json"
	models "go_pet_shop/internal/domain"
	"log/slog"
	"net/http"
)

type Users interface {
	CreateUser(user models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetAllUsers() ([]models.User, error)
	DeleteUser(id string) error
}

func CreateUser(log *slog.Logger, users Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.users.CreateUser"

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Error("Failed to decode request body", slog.String("function", fn), slog.String("error", err.Error()))
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		err := users.CreateUser(user)
		if err != nil {
			log.Error("Failed to create user", slog.String("function", fn), slog.String("error", err.Error()))
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": user.ID})
	}
}

func GetUserByEmail(log *slog.Logger, users Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.users.GetUserByEmail"
		email := r.URL.Query().Get("email")
		user, err := users.GetUserByEmail(email)
		if err != nil {
			log.Error("User not found", slog.String("function", fn), slog.String("error", err.Error()))
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

func GetAllUsers(log *slog.Logger, users Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.users.GetAllUsers"
		usersList, err := users.GetAllUsers()
		if err != nil {
			log.Error("Failed to retrieve users", slog.String("function", fn), slog.String("error", err.Error()))
			http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(usersList)
	}
}

func DeleteUser(log *slog.Logger, users Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.users.DeleteUser"
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}
		err := users.DeleteUser(id)
		if err != nil {
			log.Error("Failed to delete user", slog.String("function", fn), slog.String("error", err.Error()))
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
