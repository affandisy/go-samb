package service

import (
	"context"
	"go-samb/internal/model"
	"go-samb/internal/repository"
)

type WarehouseService interface {
	GetAll(ctx context.Context) ([]model.Warehouse, error)
}

type warehouseService struct {
	warehouseRepo repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) WarehouseService {
	return &warehouseService{warehouseRepo: repo}
}

func (s *warehouseService) GetAll(ctx context.Context) ([]model.Warehouse, error) {
	return s.warehouseRepo.GetAll(ctx)
}
