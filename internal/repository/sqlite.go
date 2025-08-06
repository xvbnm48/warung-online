package repository

import (
	"database/sql"
	"log"
	"warung-online/internal/entity"

	_ "github.com/mattn/go-sqlite3"
)

func NewSQLiteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS stocks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			price REAL NOT NULL,
			quantity INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			customer_name TEXT NOT NULL,
			order_date DATETIME NOT NULL,
			total_amount REAL NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			stock_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			price REAL NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id),
			FOREIGN KEY (stock_id) REFERENCES stocks(id)
		)`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

type sqliteOrderRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewSQLiteOrderRepository(db *sql.DB, logger *log.Logger) OrderRepository {
	return &sqliteOrderRepository{db, logger}
}

func (r *sqliteOrderRepository) CreateOrder(order entity.Order) (entity.Order, error) {
	r.logger.Println("CreateOrder repository started")
	defer r.logger.Println("CreateOrder repository finished")
	tx, err := r.db.Begin()
	if err != nil {
		return order, err
	}

	res, err := tx.Exec("INSERT INTO orders (customer_name, order_date, total_amount) VALUES (?, ?, ?)", order.CustomerName, order.OrderDate, order.TotalAmount)
	if err != nil {
		tx.Rollback()
		return order, err
	}
	orderID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return order, err
	}
	order.ID = orderID

	for i, item := range order.Items {
		res, err := tx.Exec("INSERT INTO order_items (order_id, stock_id, quantity, price) VALUES (?, ?, ?, ?)", orderID, item.StockID, item.Quantity, item.Price)
		if err != nil {
			tx.Rollback()
			return order, err
		}
		itemID, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return order, err
		}
		order.Items[i].ID = itemID
		order.Items[i].OrderID = orderID
	}

	return order, tx.Commit()
}

func (r *sqliteOrderRepository) GetOrder(id int64) (entity.Order, error) {
	r.logger.Printf("GetOrder repository started for id: %d", id)
	defer r.logger.Printf("GetOrder repository finished for id: %d", id)
	var order entity.Order
	row := r.db.QueryRow("SELECT id, customer_name, order_date, total_amount FROM orders WHERE id = ?", id)
	err := row.Scan(&order.ID, &order.CustomerName, &order.OrderDate, &order.TotalAmount)
	if err != nil {
		return order, err
	}

	rows, err := r.db.Query("SELECT id, order_id, stock_id, quantity, price FROM order_items WHERE order_id = ?", id)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.StockID, &item.Quantity, &item.Price); err != nil {
			return order, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *sqliteOrderRepository) ListOrders() ([]entity.Order, error) {
	r.logger.Println("ListOrders repository started")
	defer r.logger.Println("ListOrders repository finished")
	rows, err := r.db.Query("SELECT id, customer_name, order_date, total_amount FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.CustomerName, &order.OrderDate, &order.TotalAmount); err != nil {
			return nil, err
		}

		itemRows, err := r.db.Query("SELECT id, order_id, stock_id, quantity, price FROM order_items WHERE order_id = ?", order.ID)
		if err != nil {
			return nil, err // Return error if query fails
		}

		for itemRows.Next() {
			var item entity.OrderItem
			if err := itemRows.Scan(&item.ID, &item.OrderID, &item.StockID, &item.Quantity, &item.Price); err != nil {
				itemRows.Close() // Close rows before returning error
				return nil, err
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close() // Close rows after iterating

		orders = append(orders, order)
	}
	return orders, nil
}

func (r *sqliteOrderRepository) UpdateOrder(id int64, order entity.Order) (entity.Order, error) {
	// Complex logic: need to handle item updates, additions, deletions
	// For simplicity, this example only updates the main order fields
	_, err := r.db.Exec("UPDATE orders SET customer_name = ?, order_date = ?, total_amount = ? WHERE id = ?", order.CustomerName, order.OrderDate, order.TotalAmount, id)
	if err != nil {
		return order, err
	}
	order.ID = id
	return order, nil
}

func (r *sqliteOrderRepository) DeleteOrder(id int64) error {
	_, err := r.db.Exec("DELETE FROM orders WHERE id = ?", id)
	return err
}

type sqliteStockRepository struct {
	db     *sql.DB
	logger *log.Logger
}

func NewSQLiteStockRepository(db *sql.DB, logger *log.Logger) StockRepository {
	return &sqliteStockRepository{db, logger}
}

func (r *sqliteStockRepository) CreateStock(stock entity.Stock) (entity.Stock, error) {
	r.logger.Println("CreateStock repository started")
	defer r.logger.Println("CreateStock repository finished")
	res, err := r.db.Exec("INSERT INTO stocks (name, price, quantity) VALUES (?, ?, ?)", stock.Name, stock.Price, stock.Quantity)
	if err != nil {
		return stock, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return stock, err
	}
	stock.ID = id
	return stock, nil
}

func (r *sqliteStockRepository) GetStock(id int64) (entity.Stock, error) {
	r.logger.Printf("GetStock repository started for id: %d", id)
	defer r.logger.Printf("GetStock repository finished for id: %d", id)
	var stock entity.Stock
	row := r.db.QueryRow("SELECT id, name, price, quantity FROM stocks WHERE id = ?", id)
	err := row.Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Quantity)
	return stock, err
}

func (r *sqliteStockRepository) ListStocks() ([]entity.Stock, error) {
	r.logger.Println("ListStocks repository started")
	defer r.logger.Println("ListStocks repository finished")
	rows, err := r.db.Query("SELECT id, name, price, quantity FROM stocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []entity.Stock
	for rows.Next() {
		var stock entity.Stock
		if err := rows.Scan(&stock.ID, &stock.Name, &stock.Price, &stock.Quantity); err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (r *sqliteStockRepository) UpdateStock(id int64, stock entity.Stock) (entity.Stock, error) {
	_, err := r.db.Exec("UPDATE stocks SET name = ?, price = ?, quantity = ? WHERE id = ?", stock.Name, stock.Price, stock.Quantity, id)
	if err != nil {
		return stock, err
	}
	stock.ID = id
	return stock, nil
}

func (r *sqliteStockRepository) DeleteStock(id int64) error {
	_, err := r.db.Exec("DELETE FROM stocks WHERE id = ?", id)
	return err
}
