package model

import "time"

type TransactionOutHeader struct {
	TrxOutPK      int       `json:"trx_out_pk"`
	TrxOutNo      string    `json:"trx_out_no"`
	WhsIdf        int       `json:"whs_idf"`
	TrxOutDate    string    `json:"trx_out_date"`
	TrxOutCustIdf int       `json:"trx_out_cust_idf"`
	TrxOutNotes   string    `json:"trx_out_notes"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransactionOutDetail struct {
	TrxOutDPK         int       `json:"trx_out_d_pk"`
	TrxOutIdf         int       `json:"trx_out_idf"`
	TrxOutDProductIdf int       `json:"trx_out_d_product_idf"`
	TrxOutDQtyDus     int       `json:"trx_out_d_qty_dus"`
	TrxOutDQtyPcs     int       `json:"trx_out_d_qty_pcs"`
	CreatedAt         time.Time `json:"created_at"`
}

type TransactionOut struct {
	Header  TransactionOutHeader   `json:"header"`
	Details []TransactionOutDetail `json:"details"`
}

type TransactionOutList struct {
	Header        TransactionOutHeader `json:"header"`
	WarehouseName string               `json:"warehouse_name"`
	CustomerName  string               `json:"customer_name"`
}

type TransactionOutDetailView struct {
	Detail      TransactionOutDetail `json:"detail"`
	ProductName string               `json:"product_name"`
}

type CreateTransactionOutRequest struct {
	TrxOutNo      string                              `json:"trx_out_no"`
	WhsIdf        int                                 `json:"whs_idf"`
	TrxOutDate    string                              `json:"trx_out_date"`
	TrxOutCustIdf int                                 `json:"trx_out_cust_idf"`
	TrxOutNotes   string                              `json:"trx_out_notes"`
	Details       []CreateTransactionOutDetailRequest `json:"details"`
}

type CreateTransactionOutDetailRequest struct {
	TrxOutDProductIdf int `json:"trx_out_d_product_idf"`
	TrxOutDQtyDus     int `json:"trx_out_d_qty_dus"`
	TrxOutDQtyPcs     int `json:"trx_out_d_qty_pcs"`
}
