package postgres

import (
	"context"
	"fmt"
	models "go_pet_shop/internal/domain"
)

func (s *Storage) PlaceOrder(userEmail string, items []models.OrderItem) (orderID int, err error) {
	const fn = "storage.postgres.transaction.PlaceOrder"
	ctx := context.Background()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	for _, item := range items {
		res, execErr := tx.Exec(ctx, `
			UPDATE products 
			SET stock = stock - $1
			WHERE id = $2 AND stock >= $1
			`, item.Quantity, item.ProductID)
		if execErr != nil {
			return 0, fmt.Errorf("%s: update stock failed: %w", fn, execErr)
		}
		if res.RowsAffected() == 0 {
			return 0, fmt.Errorf("%s: insufficient stock for product ID %d", fn, item.ProductID)
		}
	}

	var userID int
	err = tx.QueryRow(ctx, `
		SELECT id FROM users WHERE email = $1`, userEmail).
		Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to find user ID: %w", fn, err)
	}
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, total_price, created_at)
		VALUES ($1, 0, NOW())
		RETURNING id
		`, userID).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("%s: insert order failed: %w", fn, err)
	}

	var total float64
	for _, item := range items {
		var price float64
		err = tx.QueryRow(ctx, `
			SELECT price FROM products WHERE id = $1`, item.ProductID).Scan(&price)
		if err != nil {
			return 0, fmt.Errorf("%s: get price failed for product %d: %w", fn, item.ProductID, err)
		}
		_, execErr := tx.Exec(ctx, `
		INSERT INTO order_items (order_id, product_id, quantity)
		VALUES ($1, $2, $3)
		`, orderID, item.ProductID, item.Quantity)
		if execErr != nil {
			return 0, fmt.Errorf("%s: insert order item failed: %w", fn, execErr)
		}
		total += float64(item.Quantity) * price
	}

	_, err = tx.Exec(ctx, `
		UPDATE orders SET total_price = $1 WHERE id = $2
		`, total, orderID)
	if err != nil {
		return 0, fmt.Errorf("%s: update order total failed: %w", fn, err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO transactions (order_id, amount, status, created_at)
		VALUES ($1, $2, 'paid', NOW())
		`, orderID, total)
	if err != nil {
		return 0, fmt.Errorf("%s: insert transction failed: %w", fn, err)
	}
	return orderID, nil
}
