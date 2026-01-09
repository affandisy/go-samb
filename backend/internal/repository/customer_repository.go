package repository

import (
	"context"
	"database/sql"
	"go-samb/internal/model"
)

type CustomerRepository interface {
	GetAll(ctx context.Context) ([]model.Customer, error)
	GetByID(ctx context.Context, id int) (*model.Customer, error)
}

type customerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *customerRepository {
	return &customerRepository{
		DB: db,
	}
}

func (r *customerRepository) GetAll(ctx context.Context) ([]model.Customer, error) {
	query := "SELECT customer_pk, customer_name, created_at FROM master_customer ORDER BY customer_name"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(&c.CustomerPK, &c.CustomerName, &c.CreatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	return customers, err

}

func (r *customerRepository) GetByID(ctx context.Context, id int) (*model.Customer, error) {
	query := "SELECT customer_pk, customer_name, created_at FROM master_customer WHERE customer_pk = $1"

	var c model.Customer
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&c.CustomerPK, &c.CustomerName, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
