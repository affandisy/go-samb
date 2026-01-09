package service

import (
	"context"
	"go-samb/internal/model"
	"go-samb/internal/repository"
)

type StockService interface {
	GetStockReport(ctx context.Context) ([]model.StockReport, error)
}

type stockService struct {
	stockRepo repository.StockRepository
}

func NewStockService(repo repository.StockRepository) StockService {
	return &stockService{stockRepo: repo}
}

func (s *stockService) GetStockReport(ctx context.Context) ([]model.StockReport, error) {
	return s.stockRepo.GetStockReport(ctx)
}
