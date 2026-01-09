package repository

import (
	"context"
	"database/sql"
	"go-samb/internal/model"
)

type SupplierRepository interface {
	GetAll(ctx context.Context) ([]model.Supplier, error)
	GetByID(ctx context.Context, id int) (*model.Supplier, error)
}

type supplierRepository struct {
	DB *sql.DB
}

func NewSupplierRepository(db *sql.DB) *supplierRepository {
	return &supplierRepository{
		DB: db,
	}
}

func (r *supplierRepository) GetAll(ctx context.Context) ([]model.Supplier, error) {
	query := "SELECT supplier_pk, supplier_name, created_at FROM master_supplier ORDER BY supplier_name"
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []model.Supplier
	for rows.Next() {
		var s model.Supplier
		if err := rows.Scan(&s.SupplierPK, &s.SupplierName, &s.CreatedAt); err != nil {
			return nil, err
		}
		suppliers = append(suppliers, s)
	}

	return suppliers, err

}

func (r *supplierRepository) GetByID(ctx context.Context, id int) (*model.Supplier, error) {
	query := "SELECT supplier_pk, supplier_name, created_at FROM master_supplier WHERE supplier_pk = $1"

	var s model.Supplier
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&s.SupplierPK, &s.SupplierName, &s.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
