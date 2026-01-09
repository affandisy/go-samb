package main

import (
	"fmt"
	"go-samb/internal/handler"
	"go-samb/internal/repository"
	"go-samb/internal/router"
	"go-samb/internal/service"
	"go-samb/pkg/config"
	"go-samb/pkg/database"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("Main Go Init")

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("Gagal load timezone jakarta")
	} else {
		time.Local = loc
	}

	// dbURL := os.Getenv("DATABASE_URL")
	// if dbURL == "" {
	// 	log.Fatal("environment variable is required")
	// }

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	// Initialize Repositories
	supplierRepo := repository.NewSupplierRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	productRepo := repository.NewProductRepository(db)
	warehouseRepo := repository.NewWarehouseRepository(db)
	transactionInRepo := repository.NewTransactionInRepository(db)
	transactionOutRepo := repository.NewTransactionOutRepository(db)
	stockRepo := repository.NewStockRepository(db)

	// Init Service
	supplierSvc := service.NewSupplierService(supplierRepo)
	customerSvc := service.NewCustomerService(customerRepo)
	productSvc := service.NewProductService(productRepo)
	warehouseSvc := service.NewWarehouseService(warehouseRepo)
	transactionInSvc := service.NewTransactionInService(transactionInRepo)
	transactionOutSvc := service.NewTransactionOutService(transactionOutRepo, stockRepo)
	stockSvc := service.NewStockService(stockRepo)

	// Init Handler
	masterHandler := handler.NewMasterHandler(supplierSvc, customerSvc, productSvc, warehouseSvc)
	transactionInHandler := handler.NewTransactionInHandler(transactionInSvc)
	transactionOutHandler := handler.NewTransactionOutHandler(transactionOutSvc)
	stockHandler := handler.NewStockHandler(stockSvc)

	e := echo.New()

	router.SetupRouter(e, &router.RouterConfig{
		MasterHandler:         masterHandler,
		TransactionInHandler:  transactionInHandler,
		TransactionOutHandler: transactionOutHandler,
		StockHandler:          stockHandler,
	})

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatal("failed to start server")
	}
}
