package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-samb/internal/model"
	"strings"
)

type StockRepository interface {
	GetStockReport(ctx context.Context) ([]model.StockReport, error)
	GetCurrentStock(ctx context.Context, whsID int, productID int) (int, int, error)
	GetCurrentStockBatch(ctx context.Context, whsID int, productIDs []int) (map[int]model.StockData, error)
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
	// query := `SELECT warehouse, product, qty_dus, qty_pcs FROM vw_stock_report WHERE qty_dus > 0 OR qty_pcs > 0 ORDER BY warehouse, product`
	query := `SELECT warehouse, product, qty_dus, qty_pcs 
	          FROM vw_stock_report 
	          ORDER BY warehouse, product`
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
		WITH stock_in AS (
			SELECT 
				COALESCE(SUM(tid.trx_in_d_qty_dus), 0) as total_dus,
				COALESCE(SUM(tid.trx_in_d_qty_pcs), 0) as total_pcs
			FROM trx_in_detail tid
			JOIN trx_in_header tih ON tid.trx_in_idf = tih.trx_in_pk
			WHERE tih.whs_idf = $1 AND tid.trx_in_d_product_idf = $2
		),
		stock_out AS (
			SELECT 
				COALESCE(SUM(tod.trx_out_d_qty_dus), 0) as total_dus,
				COALESCE(SUM(tod.trx_out_d_qty_pcs), 0) as total_pcs
			FROM trx_out_detail tod
			JOIN trx_out_header toh ON tod.trx_out_idf = toh.trx_out_pk
			WHERE toh.whs_idf = $1 AND tod.trx_out_d_product_idf = $2
		)
		SELECT 
			(si.total_dus - so.total_dus) as current_dus,
			(si.total_pcs - so.total_pcs) as current_pcs
		FROM stock_in si, stock_out so
	`

	var currentDus, currentPcs int
	err := r.DB.QueryRowContext(ctx, query, whsID, productID).Scan(&currentDus, &currentPcs)
	if err != nil {
		return 0, 0, err
	}

	return currentDus, currentPcs, nil
}

func (r *stockRepository) GetCurrentStockBatch(ctx context.Context, whsID int, productIDs []int) (map[int]model.StockData, error) {
	if len(productIDs) == 0 {
		return make(map[int]model.StockData), nil
	}

	// Build placeholders for IN clause
	placeholders := make([]string, len(productIDs))
	args := make([]interface{}, len(productIDs)+1)
	args[0] = whsID

	for i, id := range productIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = id
	}

	query := fmt.Sprintf(`
		WITH stock_in AS (
			SELECT 
				tid.trx_in_d_product_idf as product_id,
				COALESCE(SUM(tid.trx_in_d_qty_dus), 0) as total_dus,
				COALESCE(SUM(tid.trx_in_d_qty_pcs), 0) as total_pcs
			FROM trx_in_detail tid
			JOIN trx_in_header tih ON tid.trx_in_idf = tih.trx_in_pk
			WHERE tih.whs_idf = $1 AND tid.trx_in_d_product_idf IN (%s)
			GROUP BY tid.trx_in_d_product_idf
		),
		stock_out AS (
			SELECT 
				tod.trx_out_d_product_idf as product_id,
				COALESCE(SUM(tod.trx_out_d_qty_dus), 0) as total_dus,
				COALESCE(SUM(tod.trx_out_d_qty_pcs), 0) as total_pcs
			FROM trx_out_detail tod
			JOIN trx_out_header toh ON tod.trx_out_idf = toh.trx_out_pk
			WHERE toh.whs_idf = $1 AND tod.trx_out_d_product_idf IN (%s)
			GROUP BY tod.trx_out_d_product_idf
		)
		SELECT 
			COALESCE(si.product_id, so.product_id) as product_id,
			COALESCE(si.total_dus, 0) - COALESCE(so.total_dus, 0) as current_dus,
			COALESCE(si.total_pcs, 0) - COALESCE(so.total_pcs, 0) as current_pcs
		FROM stock_in si
		FULL OUTER JOIN stock_out so ON si.product_id = so.product_id
	`, strings.Join(placeholders, ","), strings.Join(placeholders, ","))

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]model.StockData)
	for rows.Next() {
		var productID, dus, pcs int
		if err := rows.Scan(&productID, &dus, &pcs); err != nil {
			return nil, err
		}
		result[productID] = model.StockData{Dus: dus, Pcs: pcs}
	}

	return result, nil
}
