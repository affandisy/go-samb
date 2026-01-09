package handler

import (
	"go-samb/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MasterHandler struct {
	supplierService  service.SupplierService
	customerService  service.CustomerService
	productService   service.ProductService
	warehouseService service.WarehouseService
}

func NewMasterHandler(supplierSvc service.SupplierService,
	customerSvc service.CustomerService,
	productSvc service.ProductService,
	warehouseSvc service.WarehouseService) *MasterHandler {
	return &MasterHandler{
		supplierService:  supplierSvc,
		customerService:  customerSvc,
		productService:   productSvc,
		warehouseService: warehouseSvc,
	}
}

func (h *MasterHandler) GetSuppliers(c echo.Context) error {
	ctx := c.Request().Context()
	suppliers, err := h.supplierService.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, suppliers)
}

func (h *MasterHandler) GetCustomers(c echo.Context) error {
	ctx := c.Request().Context()
	customers, err := h.customerService.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, customers)
}

func (h *MasterHandler) GetProducts(c echo.Context) error {
	ctx := c.Request().Context()
	products, err := h.productService.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

func (h *MasterHandler) GetWarehouses(c echo.Context) error {
	ctx := c.Request().Context()
	warehouses, err := h.warehouseService.GetAll(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, warehouses)
}
