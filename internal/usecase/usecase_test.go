package usecase

import (
	"errors"
	"log"
	"os"
	"testing"
	"warung-online/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository adalah mock untuk OrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(order entity.Order) (entity.Order, error) {
	args := m.Called(order)
	return args.Get(0).(entity.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrder(id int64) (entity.Order, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Order), args.Error(1)
}

func (m *MockOrderRepository) ListOrders() ([]entity.Order, error) {
	args := m.Called()
	return args.Get(0).([]entity.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrder(id int64, order entity.Order) (entity.Order, error) {
	args := m.Called(id, order)
	return args.Get(0).(entity.Order), args.Error(1)
}

func (m *MockOrderRepository) DeleteOrder(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockStockRepository adalah mock untuk StockRepository
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) CreateStock(stock entity.Stock) (entity.Stock, error) {
	args := m.Called(stock)
	return args.Get(0).(entity.Stock), args.Error(1)
}

func (m *MockStockRepository) GetStock(id int64) (entity.Stock, error) {
	args := m.Called(id)
	return args.Get(0).(entity.Stock), args.Error(1)
}

func (m *MockStockRepository) ListStocks() ([]entity.Stock, error) {
	args := m.Called()
	return args.Get(0).([]entity.Stock), args.Error(1)
}

func (m *MockStockRepository) UpdateStock(id int64, stock entity.Stock) (entity.Stock, error) {
	args := m.Called(id, stock)
	return args.Get(0).(entity.Stock), args.Error(1)
}

func (m *MockStockRepository) DeleteStock(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func newTestLogger() *log.Logger {
	return log.New(os.Stdout, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func TestCreateOrder_Success(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockStockRepo := new(MockStockRepository)
	logger := newTestLogger()
	usecase := NewOrderUsecase(mockOrderRepo, mockStockRepo, logger)

	order := entity.Order{
		CustomerName: "Test Customer",
		Items: []entity.OrderItem{
			{StockID: 1, Quantity: 2},
		},
	}

	stock := entity.Stock{ID: 1, Name: "Test Stock", Price: 100, Quantity: 10}

	// Mocking
	mockStockRepo.On("GetStock", int64(1)).Return(stock, nil)

	// Buat mock untuk order yang dikembalikan, yang sudah memiliki TotalAmount
	// dan item yang diperbarui.
	returnedOrder := order
	returnedOrder.ID = 1 // ID ditetapkan oleh repo
	returnedOrder.TotalAmount = 200.0
	returnedOrder.Items[0].Price = 100.0
	returnedOrder.Items[0].ID = 1
	returnedOrder.Items[0].OrderID = 1

	mockOrderRepo.On("CreateOrder", mock.AnythingOfType("entity.Order")).Return(returnedOrder, nil)
	mockStockRepo.On("UpdateStock", int64(1), mock.AnythingOfType("entity.Stock")).Return(entity.Stock{}, nil)

	// Execute
	createdOrder, err := usecase.CreateOrder(order)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200.0, createdOrder.TotalAmount)
	mockStockRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func TestCreateOrder_StockNotFound(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockStockRepo := new(MockStockRepository)
	logger := newTestLogger()
	usecase := NewOrderUsecase(mockOrderRepo, mockStockRepo, logger)

	order := entity.Order{
		CustomerName: "Test Customer",
		Items: []entity.OrderItem{
			{StockID: 1, Quantity: 2},
		},
	}

	// Mocking
	mockStockRepo.On("GetStock", int64(1)).Return(entity.Stock{}, errors.New("stock not found"))

	// Execute
	_, err := usecase.CreateOrder(order)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "stock not found", err.Error())
	mockStockRepo.AssertExpectations(t)
}

func TestCreateOrder_NotEnoughStock(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockStockRepo := new(MockStockRepository)
	logger := newTestLogger()
	usecase := NewOrderUsecase(mockOrderRepo, mockStockRepo, logger)

	order := entity.Order{
		CustomerName: "Test Customer",
		Items: []entity.OrderItem{
			{StockID: 1, Quantity: 10},
		},
	}

	stock := entity.Stock{ID: 1, Name: "Test Stock", Price: 100, Quantity: 5}

	// Mocking
	mockStockRepo.On("GetStock", int64(1)).Return(stock, nil)

	// Execute
	_, err := usecase.CreateOrder(order)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "not enough stock", err.Error())
	mockStockRepo.AssertExpectations(t)
}
