package handler

import (
	"go-samb/internal/model"
	"go-samb/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TransactionInHandler struct {
	trxInService service.TransactionInService
}

func NewTransactionInHandler(trxInSvc service.TransactionInService) *TransactionInHandler {
	return &TransactionInHandler{trxInService: trxInSvc}
}

func (h *TransactionInHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req model.CreateTransactionInRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Convert request to domain model
	trx := &model.TransactionIn{
		Header: model.TransactionInHeader{
			TrxInNo:      req.TrxInNo,
			WhsIdf:       req.WhsIdf,
			TrxInDate:    req.TrxInDate,
			TrxInSuppIdf: req.TrxInSuppIdf,
			TrxInNotes:   req.TrxInNotes,
		},
		Details: make([]model.TransactionInDetail, len(req.Details)),
	}

	for i, detail := range req.Details {
		trx.Details[i] = model.TransactionInDetail{
			TrxInDProductIdf: detail.TrxInDProductIdf,
			TrxInDQtyDus:     detail.TrxInDQtyDus,
			TrxInDQtyPcs:     detail.TrxInDQtyPcs,
		}
	}

	id, err := h.trxInService.Create(ctx, trx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "transaction created",
		"id":      id,
	})
}

func (h *TransactionInHandler) GetList(c echo.Context) error {
	ctx := c.Request().Context()
	transactions, err := h.trxInService.GetList(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h *TransactionInHandler) GetDetail(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	header, details, err := h.trxInService.GetDetail(ctx, idInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "transaction not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"header":  header,
		"details": details,
	})
}
