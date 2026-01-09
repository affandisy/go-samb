package repository

import (
	"context"
	"database/sql"
	"go-samb/internal/model"
)

type WarehouseRepository interface {
	GetAll(ctx context.Context) ([]model.Warehouse, error)
	GetByID(ctx context.Context, id int) (*model.Warehouse, error)
}

type warehouseRepository struct {
	DB *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *warehouseRepository {
	return &warehouseRepository{
		DB: db,
	}
}

func (r *warehouseRepository) GetAll(ctx context.Context) ([]model.Warehouse, error) {
	query := "SELECT whs_pk, whs_name, created_at FROM master_warehouse ORDER BY whs_name"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []model.Warehouse
	for rows.Next() {
		var w model.Warehouse
		if err := rows.Scan(&w.WhsPK, &w.WhsName, &w.CreatedAt); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, w)
	}

	return warehouses, err

}

func (r *warehouseRepository) GetByID(ctx context.Context, id int) (*model.Warehouse, error) {
	query := "SELECT whs_pk, whs_name, created_at FROM master_warehouse WHERE whs_pk = $1"

	var w model.Warehouse
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&w.WhsPK, &w.WhsName, &w.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &w, nil
}
