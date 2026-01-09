package repository

import (
	"context"
	"database/sql"
	"go-samb/internal/model"
)

type TransactionInRepository interface {
	Create(ctx context.Context, trx *model.TransactionIn) (int, error)
	GetList(ctx context.Context) ([]model.TransactionInList, error)
	GetDetail(ctx context.Context, id int) (*model.TransactionInHeader, []model.TransactionInDetailView, error)
}

type transactionRepository struct {
	DB *sql.DB
}

func NewTransactionInRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{
		DB: db,
	}
}

func (r *transactionRepository) Create(ctx context.Context, trx *model.TransactionIn) (int, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback()

	var headerID int
	query := `INSERT INTO trx_in_header (trx_in_no, whs_idf, trx_in_date, trx_in_supp_idf, trx_in_notes) VALUES ($1, $2, $3, $4, $5) RETURNING trx_in_pk`
	err = tx.QueryRowContext(ctx, query, trx.Header.TrxInNo, trx.Header.WhsIdf, trx.Header.TrxInDate, trx.Header.TrxInSuppIdf, trx.Header.TrxInNotes).Scan(&headerID)
	if err != nil {
		return 0, err
	}

	detailQuery := `INSERT INTO trx_in_detail (trx_in_idf, trx_in_d_product_idf, trx_in_d_qty_dus, trx_in_d_qty_pcs) VALUES ($1, $2, $3, $4)`
	for _, detail := range trx.Details {
		_, err = tx.ExecContext(ctx, detailQuery, headerID, detail.TrxInDProductIdf, detail.TrxInDQtyDus, detail.TrxInDQtyPcs)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return headerID, nil

}

func (r *transactionRepository) GetList(ctx context.Context) ([]model.TransactionInList, error) {
	query := `SELECT tih.trx_in_pk, tih.trx_in_no, tih.whs_idf, tih.trx_in_date, 
		tih.trx_in_supp_idf, tih.trx_in_notes, tih.created_at,
		w.whs_name, s.supplier_name
		FROM trx_in_header tih
		JOIN master_warehouse w ON tih.whs_idf = w.whs_pk
		JOIN master_supplier s ON tih.trx_in_supp_idf = s.supplier_pk
		ORDER BY tih.trx_in_date DESC, tih.trx_in_pk DESC`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.TransactionInList
	for rows.Next() {
		var t model.TransactionInList
		if err := rows.Scan(
			&t.Header.TrxInPK,
			&t.Header.TrxInNo,
			&t.Header.WhsIdf,
			&t.Header.TrxInDate,
			&t.Header.TrxInSuppIdf,
			&t.Header.TrxInNotes,
			&t.Header.CreatedAt,
			&t.WarehouseName,
			&t.SupplierName); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *transactionRepository) GetDetail(ctx context.Context, id int) (*model.TransactionInHeader, []model.TransactionInDetailView, error) {
	headerQuery := `SELECT trx_in_pk, trx_in_no, whs_idf, trx_in_date, trx_in_supp_idf, trx_in_notes, created_at
		FROM trx_in_header WHERE trx_in_pk = $1`

	var header model.TransactionInHeader
	err := r.DB.QueryRowContext(ctx, headerQuery, id).Scan(
		&header.TrxInPK,
		&header.TrxInNo,
		&header.WhsIdf,
		&header.TrxInDate,
		&header.TrxInSuppIdf,
		&header.TrxInNotes,
		&header.CreatedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	detailQuery := `SELECT tid.trx_in_d_pk, tid.trx_in_idf, tid.trx_in_d_product_idf, 
		tid.trx_in_d_qty_dus, tid.trx_in_d_qty_pcs, tid.created_at,
		p.product_name
		FROM trx_in_detail tid
		JOIN master_product p ON tid.trx_in_d_product_idf = p.product_pk
		WHERE tid.trx_in_idf = $1`
	rows, err := r.DB.QueryContext(ctx, detailQuery, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var details []model.TransactionInDetailView
	for rows.Next() {
		var d model.TransactionInDetailView
		if err := rows.Scan(
			&d.Detail.TrxInDPK,
			&d.Detail.TrxInIdf,
			&d.Detail.TrxInDProductIdf,
			&d.Detail.TrxInDQtyDus,
			&d.Detail.TrxInDQtyPcs,
			&d.Detail.CreatedAt,
			&d.ProductName,
		); err != nil {
			return nil, nil, err
		}
		details = append(details, d)
	}

	return &header, details, nil
}
