package api

import (
	"net/http"
	"strconv"
	"warung-online/internal/entity"
	"warung-online/internal/usecase"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase usecase.OrderUsecase
}

func NewOrderHandler(usecase usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: usecase}
}

type CreateOrderRequest struct {
	CustomerName string             `json:"customer_name"`
	Items        []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	StockID  int64 `json:"stock_id"`
	Quantity int   `json:"quantity"`
}

// CreateOrder godoc
// @Summary Create an order
// @Description Create a new order with multiple items
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body CreateOrderRequest true "Order Request"
// @Success 201 {object} entity.Order
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	order := entity.Order{
		CustomerName: req.CustomerName,
	}
	for _, itemReq := range req.Items {
		order.Items = append(order.Items, entity.OrderItem{
			StockID:  itemReq.StockID,
			Quantity: itemReq.Quantity,
		})
	}

	createdOrder, err := h.usecase.CreateOrder(order)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccess(c, http.StatusCreated, createdOrder)
}

// GetOrder godoc
// @Summary Get an order
// @Description Get an order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} entity.Order
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	order, err := h.usecase.GetOrder(id)
	if err != nil {
		SendError(c, http.StatusNotFound, "Order not found")
		return
	}
	SendSuccess(c, http.StatusOK, order)
}

// ListOrders godoc
// @Summary List orders
// @Description Get a list of orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Order
// @Router /orders [get]
func (h *OrderHandler) ListOrders(c *gin.Context) {
	orders, err := h.usecase.ListOrders()
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccess(c, http.StatusOK, orders)
}

// UpdateOrder godoc
// @Summary Update an order
// @Description Update an order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Param order body entity.Order true "Order"
// @Success 200 {object} entity.Order
// @Router /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var order entity.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	updatedOrder, err := h.usecase.UpdateOrder(id, order)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccess(c, http.StatusOK, updatedOrder)
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete an order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 204
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.usecase.DeleteOrder(id); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccess(c, http.StatusNoContent, nil)
}

type StockHandler struct {
	usecase usecase.StockUsecase
}

func NewStockHandler(usecase usecase.StockUsecase) *StockHandler {
	return &StockHandler{usecase: usecase}
}

// CreateStock godoc
// @Summary Create a stock
// @Description Create a new stock
// @Tags stocks
// @Accept  json
// @Produce  json
// @Param stock body entity.Stock true "Stock"
// @Success 201 {object} entity.Stock
// @Router /stocks [post]
func (h *StockHandler) CreateStock(c *gin.Context) {
	var stock entity.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	createdStock, err := h.usecase.CreateStock(stock)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccess(c, http.StatusCreated, createdStock)
}

// GetStock godoc
// @Summary Get a stock
// @Description Get a stock by ID
// @Tags stocks
// @Accept  json
// @Produce  json
// @Param id path int true "Stock ID"
// @Success 200 {object} entity.Stock
// @Router /stocks/{id} [get]
func (h *StockHandler) GetStock(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	stock, err := h.usecase.GetStock(id)
	if err != nil {
		SendError(c, http.StatusNotFound, "Stock not found")
		return
	}
	SendSuccess(c, http.StatusOK, stock)
}

// ListStocks godoc
// @Summary List stocks
// @Description Get a list of stocks
// @Tags stocks
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Stock
// @Router /stocks [get]
func (h *StockHandler) ListStocks(c *gin.Context) {
	stocks, err := h.usecase.ListStocks()
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccess(c, http.StatusOK, stocks)
}

// UpdateStock godoc
// @Summary Update a stock
// @Description Update a stock by ID
// @Tags stocks
// @Accept  json
// @Produce  json
// @Param id path int true "Stock ID"
// @Param stock body entity.Stock true "Stock"
// @Success 200 {object} entity.Stock
// @Router /stocks/{id} [put]
func (h *StockHandler) UpdateStock(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	var stock entity.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	updatedStock, err := h.usecase.UpdateStock(id, stock)
	if err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccess(c, http.StatusOK, updatedStock)
}

// DeleteStock godoc
// @Summary Delete a stock
// @Description Delete a stock by ID
// @Tags stocks
// @Accept  json
// @Produce  json
// @Param id path int true "Stock ID"
// @Success 204
// @Router /stocks/{id} [delete]
func (h *StockHandler) DeleteStock(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		SendError(c, http.StatusBadRequest, "Invalid ID")
		return
	}
	if err := h.usecase.DeleteStock(id); err != nil {
		SendError(c, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccess(c, http.StatusNoContent, nil)
}
