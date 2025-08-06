package api

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(orderHandler *OrderHandler, stockHandler *StockHandler) *gin.Engine {
	r := gin.Default()

	order := r.Group("/orders")
	{
		order.POST("/", orderHandler.CreateOrder)
		order.GET("/:id", orderHandler.GetOrder)
		order.GET("/", orderHandler.ListOrders)
		order.PUT("/:id", orderHandler.UpdateOrder)
		order.DELETE("/:id", orderHandler.DeleteOrder)
	}

	stock := r.Group("/stocks")
	{
		stock.POST("/", stockHandler.CreateStock)
		stock.GET("/:id", stockHandler.GetStock)
		stock.GET("/", stockHandler.ListStocks)
		stock.PUT("/:id", stockHandler.UpdateStock)
		stock.DELETE("/:id", stockHandler.DeleteStock)
	}

	return r
}
