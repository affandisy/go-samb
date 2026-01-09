package service

import (
	"context"
	"go-samb/internal/model"
	"go-samb/internal/repository"
)

type SupplierService interface {
	GetAll(ctx context.Context) ([]model.Supplier, error)
}

type supplierService struct {
	supplierRepo repository.SupplierRepository
}

func NewSupplierService(repo repository.SupplierRepository) SupplierService {
	return &supplierService{
		supplierRepo: repo,
	}
}

func (s *supplierService) GetAll(ctx context.Context) ([]model.Supplier, error) {
	return s.supplierRepo.GetAll(ctx)
}
