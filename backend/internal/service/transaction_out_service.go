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

	for _, detail := range trx.Details {

		if detail.TrxOutDQtyDus < 0 || detail.TrxOutDQtyPcs < 0 {
			return 0, errors.New("the quantity of goods out cannot be negative")
		}

		currentDus, currentPcs, err := s.stockRepo.GetCurrentStock(ctx, trx.Header.WhsIdf, detail.TrxOutDProductIdf)
		if err != nil {
			return 0, errors.New("failed to check product stock ID")
		}

		if currentDus < detail.TrxOutDQtyDus {
			return 0, fmt.Errorf("Insufficient Box stock for product ID %d. Available: %d, Requested: %d", detail.TrxOutDProductIdf, currentDus, detail.TrxOutDQtyDus)
		}

		if currentPcs < detail.TrxOutDQtyPcs {
			return 0, fmt.Errorf("Insufficient stock of Pcs for product ID %d. Available: %d, Requested: %d", detail.TrxOutDProductIdf, currentPcs, detail.TrxOutDQtyPcs)
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
