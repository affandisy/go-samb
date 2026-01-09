package service

import (
	"context"
	"go-samb/internal/model"
	"go-samb/internal/repository"
)

type CustomerService interface {
	GetAll(ctx context.Context) ([]model.Customer, error)
}

type customerService struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{customerRepo: repo}
}

func (s *customerService) GetAll(ctx context.Context) ([]model.Customer, error) {
	return s.customerRepo.GetAll(ctx)
}
