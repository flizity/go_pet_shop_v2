package postgres

import (
	"context"
	"fmt"
	models "go_pet_shop/internal/domain"
)

func (s *Storage) GetUserOrderHistory(email string) ([]models.OrderDetail, error) {
	const fn = "storage.postgres.history.GetUserOrderHistory"

	rows, err := s.db.Query(context.Background(), `
SELECT orders.id, products.name, order_items.quantity, products.price, COALESCE(transactions.status, ''), orders.created_at
FROM orders
JOIN users ON orders.user_id = users.id
JOIN order_items ON orders.id = order_items.order_id
JOIN products ON order_items.product_id = products.id
LEFT JOIN transactions ON orders.id = transactions.order_id
WHERE users.email = $1
ORDER BY orders.created_at DESC
`, email)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", fn, err)
	}
	defer rows.Close()

	var history []models.OrderDetail
	for rows.Next() {
		var od models.OrderDetail
		err := rows.Scan(&od.OrderID, &od.ProductName, &od.Quantity, &od.Price, &od.TransactionStatus, &od.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: rows scan failed: %w", fn, err)
		}
		history = append(history, od)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", fn, err)
	}
	return history, nil
}

func (s *Storage) GetPopularProducts() ([]models.PopularProduct, error) {
	const fn = "storage.postgres.history.GetPopularProducts"

	rows, err := s.db.Query(context.Background(), `
SELECT products.id, products.name, SUM(order_items.quantity)
FROM order_items
JOIN products ON order_items.product_id = products.id
GROUP BY products.id, products.name
ORDER BY SUM(order_items.quantity) DESC
`)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", fn, err)
	}
	defer rows.Close()
	var products []models.PopularProduct
	for rows.Next() {
		var pp models.PopularProduct
		err := rows.Scan(&pp.ProductID, &pp.ProductName, &pp.TotalSold)
		if err != nil {
			return nil, fmt.Errorf("%s: rows scan failed: %w", fn, err)
		}
		products = append(products, pp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration error: %w", fn, err)
	}
	return products, nil
}
