package repository

import (
	"context"
	"database/sql"
	"go-samb/internal/model"
)

type TransactionOutRepository interface {
	Create(ctx context.Context, trx *model.TransactionOut) (int, error)
	GetList(ctx context.Context) ([]model.TransactionOutList, error)
	GetDetail(ctx context.Context, id int) (*model.TransactionOutHeader, []model.TransactionOutDetailView, error)
}

type transactionOutRepository struct {
	DB *sql.DB
}

func NewTransactionOutRepository(db *sql.DB) *transactionOutRepository {
	return &transactionOutRepository{
		DB: db,
	}
}

func (r *transactionOutRepository) Create(ctx context.Context, trx *model.TransactionOut) (int, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var headerID int
	query := `INSERT INTO trx_out_header (trx_out_no, whs_idf, trx_out_date, trx_out_cust_idf, trx_out_notes) VALUES ($1, $2, $3, $4, $5) RETURNING trx_out_pk`
	err = tx.QueryRowContext(ctx, query,
		trx.Header.TrxOutNo,
		trx.Header.WhsIdf,
		trx.Header.TrxOutDate,
		trx.Header.TrxOutCustIdf,
		trx.Header.TrxOutNotes,
	).Scan(&headerID)
	if err != nil {
		return 0, err
	}

	detailQuery := `INSERT INTO trx_out_detail (trx_out_idf, trx_out_d_product_idf, trx_out_d_qty_dus, trx_out_d_qty_pcs) VALUES ($1, $2, $3, $4)`
	for _, detail := range trx.Details {
		_, err = tx.ExecContext(ctx, detailQuery, headerID,
			detail.TrxOutDProductIdf,
			detail.TrxOutDQtyDus,
			detail.TrxOutDQtyPcs,
		)
		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return headerID, nil
}

func (r *transactionOutRepository) GetList(ctx context.Context) ([]model.TransactionOutList, error) {
	query := `SELECT toh.trx_out_pk, toh.trx_out_no, toh.whs_idf, toh.trx_out_date, 
		toh.trx_out_cust_idf, toh.trx_out_notes, toh.created_at,
		w.whs_name, c.customer_name
		FROM trx_out_header toh
		JOIN master_warehouse w ON toh.whs_idf = w.whs_pk
		JOIN master_customer c ON toh.trx_out_cust_idf = c.customer_pk
		ORDER BY toh.trx_out_date DESC, toh.trx_out_pk DESC`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.TransactionOutList
	for rows.Next() {
		var t model.TransactionOutList
		if err := rows.Scan(
			&t.Header.TrxOutPK,
			&t.Header.TrxOutNo,
			&t.Header.WhsIdf,
			&t.Header.TrxOutDate,
			&t.Header.TrxOutCustIdf,
			&t.Header.TrxOutNotes,
			&t.Header.CreatedAt,
			&t.WarehouseName,
			&t.CustomerName,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *transactionOutRepository) GetDetail(ctx context.Context, id int) (*model.TransactionOutHeader, []model.TransactionOutDetailView, error) {
	headerQuery := `SELECT trx_out_pk, trx_out_no, whs_idf, trx_out_date, trx_out_cust_idf, trx_out_notes, created_at
		FROM trx_out_header WHERE trx_out_pk = $1`
	var header model.TransactionOutHeader
	err := r.DB.QueryRowContext(ctx, headerQuery, id).Scan(
		&header.TrxOutPK,
		&header.TrxOutNo,
		&header.WhsIdf,
		&header.TrxOutDate,
		&header.TrxOutCustIdf,
		&header.TrxOutNotes,
		&header.CreatedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	detailQuery := `SELECT tod.trx_out_d_pk, tod.trx_out_idf, tod.trx_out_d_product_idf, 
		tod.trx_out_d_qty_dus, tod.trx_out_d_qty_pcs, tod.created_at,
		p.product_name
		FROM trx_out_detail tod
		JOIN master_product p ON tod.trx_out_d_product_idf = p.product_pk
		WHERE tod.trx_out_idf = $1`
	rows, err := r.DB.QueryContext(ctx, detailQuery, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var details []model.TransactionOutDetailView
	for rows.Next() {
		var d model.TransactionOutDetailView
		if err := rows.Scan(
			&d.Detail.TrxOutDPK,
			&d.Detail.TrxOutIdf,
			&d.Detail.TrxOutDProductIdf,
			&d.Detail.TrxOutDQtyDus,
			&d.Detail.TrxOutDQtyPcs,
			&d.Detail.CreatedAt,
			&d.ProductName,
		); err != nil {
			return nil, nil, err
		}
		details = append(details, d)
	}

	return &header, details, nil
}
