package handler

import (
	"go-samb/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StockHandler struct {
	stockService service.StockService
}

func NewStockHandler(stockSvc service.StockService) *StockHandler {
	return &StockHandler{stockService: stockSvc}
}

func (h *StockHandler) GetStockReport(c echo.Context) error {
	ctx := c.Request().Context()
	stocks, err := h.stockService.GetStockReport(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, stocks)
}
