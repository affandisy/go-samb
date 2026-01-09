package service

import (
	"context"
	"errors"
	"fmt"
	"go-samb/internal/model"
	"go-samb/internal/repository"
)

type TransactionOutService interface {
	Create(ctx context.Context, trx *model.TransactionOut) (int, error)
	GetList(ctx context.Context) ([]model.TransactionOutList, error)
	GetDetail(ctx context.Context, id int) (*model.TransactionOutHeader, []model.TransactionOutDetailView, error)
}

type transactionOutService struct {
	trxOutRepo repository.TransactionOutRepository
	stockRepo  repository.StockRepository
}

func NewTransactionOutService(repo repository.TransactionOutRepository, stockRepo repository.StockRepository) TransactionOutService {
	return &transactionOutService{
		trxOutRepo: repo,
		stockRepo:  stockRepo,
	}
}

func (s *transactionOutService) Create(ctx context.Context, trx *model.TransactionOut) (int, error) {
	if trx.Header.TrxOutNo == "" {
		return 0, errors.New("transaction number required")
	}

	if trx.Header.WhsIdf == 0 {
		return 0, errors.New("warehouse id is required")
	}

	if trx.Header.TrxOutCustIdf == 0 {
		return 0, errors.New("customer is needed")
	}

	if len(trx.Details) == 0 {
		return 0, errors.New("need one detail")
	}

	productIDs := make(map[int]bool)
	productQtyMap := make(map[int]struct {
		dus int
		pcs int
	})

	for _, detail := range trx.Details {
		if detail.TrxOutDQtyDus < 0 || detail.TrxOutDQtyPcs < 0 {
			return 0, errors.New("the quantity of goods out cannot be negative")
		}

		productIDs[detail.TrxOutDProductIdf] = true
		productQtyMap[detail.TrxOutDProductIdf] = struct {
			dus int
			pcs int
		}{
			dus: detail.TrxOutDQtyDus,
			pcs: detail.TrxOutDQtyPcs,
		}
	}

	stockMap, err := s.stockRepo.GetCurrentStockBatch(ctx, trx.Header.WhsIdf, getKeys(productIDs))
	if err != nil {
		return 0, fmt.Errorf("failed to check stock: %v", err)
	}

	for productID, requestedQty := range productQtyMap {
		currentStock, exists := stockMap[productID]
		if !exists {
			return 0, fmt.Errorf("product ID %d not found in stock", productID)
		}

		if currentStock.Dus < requestedQty.dus {
			return 0, fmt.Errorf("insufficient Box stock for product ID %d. Available: %d, Requested: %d",
				productID, currentStock.Dus, requestedQty.dus)
		}

		if currentStock.Pcs < requestedQty.pcs {
			return 0, fmt.Errorf("insufficient stock of Pcs for product ID %d. Available: %d, Requested: %d",
				productID, currentStock.Pcs, requestedQty.pcs)
		}
	}

	return s.trxOutRepo.Create(ctx, trx)
}

func (s *transactionOutService) GetList(ctx context.Context) ([]model.TransactionOutList, error) {
	return s.trxOutRepo.GetList(ctx)
}

func (s *transactionOutService) GetDetail(ctx context.Context, id int) (*model.TransactionOutHeader, []model.TransactionOutDetailView, error) {
	return s.trxOutRepo.GetDetail(ctx, id)
}

func getKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
