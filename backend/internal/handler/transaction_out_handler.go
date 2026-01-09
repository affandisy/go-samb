package handler

import (
	"go-samb/internal/model"
	"go-samb/internal/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TransactionOutHandler struct {
	trxOutService service.TransactionOutService
}

func NewTransactionOutHandler(trxOutSvc service.TransactionOutService) *TransactionOutHandler {
	return &TransactionOutHandler{trxOutService: trxOutSvc}
}

func (h *TransactionOutHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	var req model.CreateTransactionOutRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Convert request to domain model
	trx := &model.TransactionOut{
		Header: model.TransactionOutHeader{
			TrxOutNo:      req.TrxOutNo,
			WhsIdf:        req.WhsIdf,
			TrxOutDate:    req.TrxOutDate,
			TrxOutCustIdf: req.TrxOutCustIdf,
			TrxOutNotes:   req.TrxOutNotes,
		},
		Details: make([]model.TransactionOutDetail, len(req.Details)),
	}

	for i, detail := range req.Details {
		trx.Details[i] = model.TransactionOutDetail{
			TrxOutDProductIdf: detail.TrxOutDProductIdf,
			TrxOutDQtyDus:     detail.TrxOutDQtyDus,
			TrxOutDQtyPcs:     detail.TrxOutDQtyPcs,
		}
	}

	id, err := h.trxOutService.Create(ctx, trx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "transaction created successfully",
		"id":      id,
	})
}

func (h *TransactionOutHandler) GetList(c echo.Context) error {
	ctx := c.Request().Context()
	transactions, err := h.trxOutService.GetList(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h *TransactionOutHandler) GetDetail(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	header, details, err := h.trxOutService.GetDetail(ctx, idInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "transaction not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"header":  header,
		"details": details,
	})
}
