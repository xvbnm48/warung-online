package entity

import "time"

// Stock represents a product in inventory.
// Ini menjelaskan stok untuk item apa, beserta harga dan kuantitasnya.
type Stock struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

// Order represents a customer's order.
// Ini adalah pesanan yang bisa memiliki nama (misalnya, nama pelanggan).
type Order struct {
	ID           int64       `json:"id"`
	CustomerName string      `json:"customer_name"`
	OrderDate    time.Time   `json:"order_date"`
	TotalAmount  float64     `json:"total_amount"`
	Items        []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

// OrderItem represents a single item within an order.
// Ini adalah penghubung antara Order dan Stock, memungkinkan satu pesanan memiliki banyak item.
type OrderItem struct {
	ID       int64   `json:"id"`
	OrderID  int64   `json:"order_id"`
	StockID  int64   `json:"stock_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"` // Harga per item saat pesanan dibuat
}
