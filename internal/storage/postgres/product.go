package postgres

import (
	"context"
	"fmt"
	models "go_pet_shop/internal/domain"
)

func (s *Storage) CreateProduct(product models.Product) error {
	const fn = "storage.postgres.CreateProduct"

	_, err := s.db.Exec(context.Background(), `
		INSERT INTO products (name, description, price, stock)
		VALUES ($1, $2, $3, $4)
	`, product.ID, product.Name, product.Price, product.Stock)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetProductByID(id int) (models.Product, error) {
	const fn = "storage.postgres.product.GetProductByID"

	var product models.Product

	err := s.db.QueryRow(context.Background(),
		`SELECT id, name, price, stock FROM products WHERE id = $1`,
		id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		return models.Product{}, fmt.Errorf("%s: %w", fn, err)
	}

	return product, nil
}

func (s *Storage) GetAllProducts() ([]models.Product, error) {
	const fn = "storage.postgres.product.GetAllProducts"

	rows, err := s.db.Query(context.Background(), `SELECT * FROM products`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return products, nil
}

func (s *Storage) UpdateProduct(p models.Product) error {
	const fn = "storage.postgres.product.UpdateProduct"

	_, err := s.db.Exec(context.Background(),
		`UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4`,
		p.Name, p.Price, p.Stock, p.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) DeleteProduct(id int) error {
	const fn = "storage.postgres.product.DeleteProduct"

	_, err := s.db.Exec(context.Background(),
		`DELETE FROM products WHERE id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
