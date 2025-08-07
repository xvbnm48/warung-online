package usecase

import (
	"errors"
	"log"
	"time"
	"warung-online/internal/entity"
	"warung-online/internal/repository"
)

type orderUsecase struct {
	orderRepo repository.OrderRepository
	stockRepo repository.StockRepository
	logger    *log.Logger
}

func NewOrderUsecase(orderRepo repository.OrderRepository, stockRepo repository.StockRepository, logger *log.Logger) OrderUsecase {
	return &orderUsecase{orderRepo, stockRepo, logger}
}

func (u *orderUsecase) CreateOrder(order entity.Order) (entity.Order, error) {
	u.logger.Println("CreateOrder usecase started")
	defer u.logger.Println("CreateOrder usecase finished")
	var totalAmount float64
	for i, item := range order.Items {
		stock, err := u.stockRepo.GetStock(item.StockID)
		if err != nil {
			return order, errors.New("stock not found")
		}

		if stock.Quantity < item.Quantity {
			return order, errors.New("not enough stock")
		}

		order.Items[i].Price = stock.Price
		totalAmount += stock.Price * float64(item.Quantity)
	}

	order.TotalAmount = totalAmount
	order.OrderDate = time.Now()

	createdOrder, err := u.orderRepo.CreateOrder(order)
	if err != nil {
		return order, err
	}

	// Kurangi stok
	for _, item := range createdOrder.Items {
		stock, err := u.stockRepo.GetStock(item.StockID)
		if err != nil {
			return createdOrder, errors.New("failed to get stock for update")
		}
		stock.Quantity -= item.Quantity
		_, err = u.stockRepo.UpdateStock(item.StockID, stock)
		if err != nil {
			// Here we have an inconsistent state. Order is created, but stock update failed.
			// A transaction would be the best way to handle this.
			// For now, returning an error is better than nothing.
			return createdOrder, errors.New("failed to update stock, inconsistent data")
		}
	}

	return createdOrder, nil
}

func (u *orderUsecase) GetOrder(id int64) (entity.Order, error) {
	u.logger.Printf("GetOrder usecase started for id: %d", id)
	defer u.logger.Printf("GetOrder usecase finished for id: %d", id)
	return u.orderRepo.GetOrder(id)
}

func (u *orderUsecase) ListOrders() ([]entity.Order, error) {
	u.logger.Println("ListOrders usecase started")
	defer u.logger.Println("ListOrders usecase finished")
	return u.orderRepo.ListOrders()
}

func (u *orderUsecase) UpdateOrder(id int64, order entity.Order) (entity.Order, error) {
	return u.orderRepo.UpdateOrder(id, order)
}

func (u *orderUsecase) DeleteOrder(id int64) error {
	return u.orderRepo.DeleteOrder(id)
}

type stockUsecase struct {
	repo   repository.StockRepository
	logger *log.Logger
}

func NewStockUsecase(repo repository.StockRepository, logger *log.Logger) StockUsecase {
	return &stockUsecase{repo, logger}
}

func (u *stockUsecase) CreateStock(stock entity.Stock) (entity.Stock, error) {
	u.logger.Println("CreateStock usecase started")
	defer u.logger.Println("CreateStock usecase finished")
	return u.repo.CreateStock(stock)
}

func (u *stockUsecase) GetStock(id int64) (entity.Stock, error) {
	u.logger.Printf("GetStock usecase started for id: %d", id)
	defer u.logger.Printf("GetStock usecase finished for id: %d", id)
	return u.repo.GetStock(id)
}

func (u *stockUsecase) ListStocks() ([]entity.Stock, error) {
	u.logger.Println("ListStocks usecase started")
	defer u.logger.Println("ListStocks usecase finished")
	return u.repo.ListStocks()
}

func (u *stockUsecase) UpdateStock(id int64, stock entity.Stock) (entity.Stock, error) {
	return u.repo.UpdateStock(id, stock)
}

func (u *stockUsecase) DeleteStock(id int64) error {
	return u.repo.DeleteStock(id)
}
