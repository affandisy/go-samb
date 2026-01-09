package service

import (
	"context"
	"go-samb/internal/model"
	"go-samb/internal/repository"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]model.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{productRepo: repo}
}

func (s *productService) GetAll(ctx context.Context) ([]model.Product, error) {
	return s.productRepo.GetAll(ctx)
}
