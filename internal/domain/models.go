package models

import "time"

type User struct {
	ID    int
	Name  string
	Email string
}

type Order struct {
	ID         int
	UserID     int
	TotalPrice float32
	CreatedAt  time.Time
}

type Product struct {
	ID    int
	Name  string
	Price float32
	Stock int
}

type Transaction struct {
	ID        int
	OrderId   int
	Amount    float32
	Status    string
	CreatedAt time.Time
}

type OrderDetail struct {
	OrderID           int       `json:"order_id"`
	ProductName       string    `json:"order_name"`
	Quantity          int       `json:"quantity"`
	Price             float32   `json:"price"`
	TransactionStatus string    `json:"transaction_status"`
	CreatedAt         time.Time `json:"created_at"`
}

type OrderItem struct {
	ID        int `json:"iD"`
	ProductID int `json:"product_id"`
	OrderID   int `json:"order_id"`
	Quantity  int `json:"quantity"`
}

type PopularProduct struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"name"`
	TotalSold   int    `json:"total_sold"`
}
