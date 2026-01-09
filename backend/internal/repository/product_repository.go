package repository

import (
	"context"
	"database/sql"
	"go-samb/internal/model"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]model.Product, error)
	GetByID(ctx context.Context, id int) (*model.Product, error)
}

type productRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{
		DB: db,
	}
}

func (r *productRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	query := "SELECT product_pk, product_name, created_at FROM master_product ORDER BY product_name"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ProductPK, &p.ProductName, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, err

}

func (r *productRepository) GetByID(ctx context.Context, id int) (*model.Product, error) {
	query := "SELECT product_pk, product_name, created_at FROM master_product WHERE product_pk = $1"

	var p model.Product
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&p.ProductPK, &p.ProductName, &p.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
