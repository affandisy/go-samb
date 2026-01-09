package repository

import (
	"context"
	"database/sql"
	"go-samb/internal/model"
)

type StockRepository interface {
	GetStockReport(ctx context.Context) ([]model.StockReport, error)
	GetCurrentStock(ctx context.Context, whsID int, productID int) (int, int, error)
}

type stockRepository struct {
	DB *sql.DB
}

func NewStockRepository(db *sql.DB) *stockRepository {
	return &stockRepository{
		DB: db,
	}
}

func (r *stockRepository) GetStockReport(ctx context.Context) ([]model.StockReport, error) {
	query := `SELECT warehouse, product, qty_dus, qty_pcs FROM vw_stock_report WHERE qty_dus > 0 OR qty_pcs > 0 ORDER BY warehouse, product`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []model.StockReport
	for rows.Next() {
		var s model.StockReport
		if err := rows.Scan(&s.Warehouse, &s.Product, &s.QtyDus, &s.QtyPcs); err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}

	return stocks, err
}

func (r *stockRepository) GetCurrentStock(ctx context.Context, whsID int, productID int) (int, int, error) {
	query := `
	SELECT 
			(COALESCE(SUM(CASE WHEN type = 'IN' THEN qty_dus ELSE 0 END), 0) - 
			COALESCE(SUM(CASE WHEN type = 'OUT' THEN qty_dus ELSE 0 END), 0)) as total_dus,
			(COALESCE(SUM(CASE WHEN type = 'IN' THEN qty_pcs ELSE 0 END), 0) - 
			COALESCE(SUM(CASE WHEN type = 'OUT' THEN qty_pcs ELSE 0 END), 0)) as total_pcs
		FROM (
			SELECT 'IN' as type, tid.trx_in_d_qty_dus as qty_dus, tid.trx_in_d_qty_pcs as qty_pcs
			FROM trx_in_detail tid
			JOIN trx_in_header tih ON tid.trx_in_idf = tih.trx_in_pk
			WHERE tih.whs_idf = $1 AND tid.trx_in_d_product_idf = $2
			
			UNION ALL
			
			SELECT 'OUT' as type, tod.trx_out_d_qty_dus as qty_dus, tod.trx_out_d_qty_pcs as qty_pcs
			FROM trx_out_detail tod
			JOIN trx_out_header toh ON tod.trx_out_idf = toh.trx_out_pk
			WHERE toh.whs_idf = $1 AND tod.trx_out_d_product_idf = $2
		) as combined_stock
		`

	var currentDus, currentPcs int
	err := r.DB.QueryRowContext(ctx, query, whsID, productID).Scan(&currentDus, &currentPcs)
	if err != nil {
		return 0, 0, err
	}

	return currentDus, currentPcs, nil
}
