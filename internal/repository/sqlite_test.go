package repository

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"
	"warung-online/internal/entity"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func newTestLogger() *log.Logger {
	return log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	assert.NoError(t, err)

	// Enable foreign key support
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	assert.NoError(t, err)

	err = InitDB(db)
	assert.NoError(t, err)

	return db
}

func TestRepositories(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	logger := newTestLogger()
	orderRepo := NewSQLiteOrderRepository(db, logger)
	stockRepo := NewSQLiteStockRepository(db, logger)

	t.Run("Stock Repository", func(t *testing.T) {
		stock := entity.Stock{Name: "Kopi", Price: 5000, Quantity: 50}
		createdStock, err := stockRepo.CreateStock(stock)
		assert.NoError(t, err)
		assert.NotZero(t, createdStock.ID)

		retrievedStock, err := stockRepo.GetStock(createdStock.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdStock.Name, retrievedStock.Name)
	})

	t.Run("Order Repository", func(t *testing.T) {
		stock, _ := stockRepo.CreateStock(entity.Stock{Name: "Test Stock", Price: 15000, Quantity: 100})
		order := entity.Order{
			CustomerName: "Test Customer",
			OrderDate:    time.Now(),
			TotalAmount:  30000,
			Items:        []entity.OrderItem{{StockID: stock.ID, Quantity: 2, Price: 15000}},
		}
		createdOrder, err := orderRepo.CreateOrder(order)
		assert.NoError(t, err)
		assert.NotZero(t, createdOrder.ID)

		retrievedOrder, err := orderRepo.GetOrder(createdOrder.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdOrder.ID, retrievedOrder.ID)
		assert.Len(t, retrievedOrder.Items, 1)

		orders, err := orderRepo.ListOrders()
		assert.NoError(t, err)
		assert.NotEmpty(t, orders)
		assert.Len(t, orders[0].Items, 1)
	})
}

func TestStockRepository(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	logger := newTestLogger()
	repo := NewSQLiteStockRepository(db, logger)

	t.Run("Create and Get Stock", func(t *testing.T) {
		stock := entity.Stock{Name: "Kopi", Price: 5000, Quantity: 50}
		createdStock, err := repo.CreateStock(stock)
		assert.NoError(t, err)
		assert.NotZero(t, createdStock.ID)

		retrievedStock, err := repo.GetStock(createdStock.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdStock.Name, retrievedStock.Name)
	})
}
