package service

import (
	"context"
	"errors"
	"go-samb/internal/model"
	"go-samb/internal/repository"
)

type TransactionInService interface {
	Create(ctx context.Context, trx *model.TransactionIn) (int, error)
	GetList(ctx context.Context) ([]model.TransactionInList, error)
	GetDetail(ctx context.Context, id int) (*model.TransactionInHeader, []model.TransactionInDetailView, error)
}

type transactionInService struct {
	trxInRepo repository.TransactionInRepository
}

func NewTransactionInService(repo repository.TransactionInRepository) TransactionInService {
	return &transactionInService{trxInRepo: repo}
}

func (s *transactionInService) Create(ctx context.Context, trx *model.TransactionIn) (int, error) {
	if trx.Header.TrxInNo == "" {
		return 0, errors.New("transaction number is required")
	}

	if trx.Header.WhsIdf == 0 {
		return 0, errors.New("warehouse id is required")
	}

	if trx.Header.TrxInSuppIdf == 0 {
		return 0, errors.New("supplier id is required")
	}

	if len(trx.Details) == 0 {
		return 0, errors.New("detail required")
	}

	for _, d := range trx.Details {
		if d.TrxInDProductIdf == 0 {
			return 0, errors.New("product cannot be empty")
		}
		if d.TrxInDQtyDus < 0 || d.TrxInDQtyPcs < 0 {
			return 0, errors.New("dus/pcs cannot be negative")
		}
	}

	return s.trxInRepo.Create(ctx, trx)

}

func (s *transactionInService) GetList(ctx context.Context) ([]model.TransactionInList, error) {
	return s.trxInRepo.GetList(ctx)
}

func (s *transactionInService) GetDetail(ctx context.Context, id int) (*model.TransactionInHeader, []model.TransactionInDetailView, error) {
	return s.trxInRepo.GetDetail(ctx, id)
}
