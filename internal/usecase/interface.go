package usecase

import "warung-online/internal/entity"

type OrderUsecase interface {
	CreateOrder(order entity.Order) (entity.Order, error)
	GetOrder(id int64) (entity.Order, error)
	ListOrders() ([]entity.Order, error)
	UpdateOrder(id int64, order entity.Order) (entity.Order, error)
	DeleteOrder(id int64) error
}

type StockUsecase interface {
	CreateStock(stock entity.Stock) (entity.Stock, error)
	GetStock(id int64) (entity.Stock, error)
	ListStocks() ([]entity.Stock, error)
	UpdateStock(id int64, stock entity.Stock) (entity.Stock, error)
	DeleteStock(id int64) error
}
