package postgres

import (
	"context"
	"fmt"
	models "go_pet_shop/internal/domain"
)

func (s *Storage) CreateUser(user models.User) error {
	const fn = "storage.postgres.user.CreateUser"

	_, err := s.db.Exec(context.Background(),
		`INSERT INTO users (name, email) VALUES ($1, $2)`,
		user.Name, user.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetUserByEmail(email string) (models.User, error) {
	const fn = "storage.postgres.User.GetUserByEmail"

	var user models.User
	err := s.db.QueryRow(context.Background(), "SELECT id, name, email FROM users WHERE LOWER(email) = LOWER($1)", email).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (s *Storage) GetAllUsers() ([]models.User, error) {
	const fn = "storage.postgres.User.GetAllUsers"

	rows, err := s.db.Query(context.Background(), "SELECT id, name, email FROM users")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return users, nil
}

func (s *Storage) DeleteUser(id string) error {
	const fn = "storage.postgres.User.DeleteUser"

	_, err := s.db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
