package postgres

import (
	"context"
	"fmt"
	models "go_pet_shop/internal/domain"
)

func (s *Storage) CreateOrder(order models.Order) (int, error) {
	const fn = "storage.postgres..order.CreateOrder"

	var orderID int
	err := s.db.QueryRow(context.Background(), `
INSERT INTO orders (user_id, total_price, created_at)
VALUES ($1, $2, $3)
RETURNING id
`, order.UserID, order.TotalPrice).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("%s: insert failed: %w", fn, err)
	}
	return orderID, nil
}

func (s *Storage) AddOrderItem(orderItem models.OrderItem) error {
	const fn = "storage.postgres.order.AddOrderItem"

	_, err := s.db.Exec(context.Background(), `
INSERT INTO order_items (product_id, order_id, quantity)
VALUES ($1, $2, $3)
`, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity)
	if err != nil {
		return fmt.Errorf("%s: insert failed: %w", fn, err)
	}
	return nil
}

func (s *Storage) GetOrderByID(id int) (models.Order, error) {
	const fn = "storage.postgres.order.GetOrderByID"

	var order models.Order
	err := s.db.QueryRow(context.Background(), `
	SELECT id, user_id, total_price, created_at
	FROM orders
	WHERE id = $1
`, id).Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.CreatedAt)
	if err != nil {
		return models.Order{}, fmt.Errorf("%s: query failed: %w", fn, err)
	}
	return order, nil
}

func (s *Storage) GetOrdersByUserEmail(email string) ([]models.Order, error) {
	const fn = "storage.postgres.order.GetOrdersByUserEmail"

	rows, err := s.db.Query(context.Background(), `
	SELECT orders.id, orders.user_id, orders.total_price, orders.created_at FROM orders
	JOIN users ON users.id = orders.user_id 
	WHERE users.email = $1`, email)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return orders, nil
}

func (s *Storage) GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error) {
	const fn = "storage.postgres.order.GetOrderItemsByOrderID"

	rows, err := s.db.Query(context.Background(), `
	SELECT id, product_id, order_id, quantity
	FROM order_items
	WHERE order_id = $1
`, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	for rows.Next() {
		var orderItem models.OrderItem
		if err := rows.Scan(&orderItem.ID, &orderItem.ProductID, &orderItem.OrderID, &orderItem.Quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orderItems = append(orderItems, orderItem)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return orderItems, nil
}
