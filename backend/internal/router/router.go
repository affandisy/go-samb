package router

import (
	"go-samb/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouterConfig struct {
	MasterHandler         *handler.MasterHandler
	TransactionInHandler  *handler.TransactionInHandler
	TransactionOutHandler *handler.TransactionOutHandler
	StockHandler          *handler.StockHandler
}

func SetupRouter(e *echo.Echo, config *RouterConfig) {
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/api")

	api.GET("/suppliers", config.MasterHandler.GetSuppliers)
	api.GET("/customers", config.MasterHandler.GetCustomers)
	api.GET("/products", config.MasterHandler.GetProducts)
	api.GET("/warehouses", config.MasterHandler.GetWarehouses)

	api.POST("/trx-in", config.TransactionInHandler.Create)
	api.GET("/trx-in", config.TransactionInHandler.GetList)
	api.GET("/trx-in/:id", config.TransactionInHandler.GetDetail)

	api.POST("/trx-out", config.TransactionOutHandler.Create)
	api.GET("/trx-out", config.TransactionOutHandler.GetList)
	api.GET("/trx-out/:id", config.TransactionOutHandler.GetDetail)

	api.GET("/stock-report", config.StockHandler.GetStockReport)
}
