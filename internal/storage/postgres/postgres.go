package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(databaseUrl string) (*Storage, error) {
	const fn = "storage.postgres.New"

	db, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	s.db.Close()
	return nil
}
