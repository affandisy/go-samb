package model

import "time"

type TransactionInHeader struct {
	TrxInPK      int       `json:"trx_in_pk"`
	TrxInNo      string    `json:"trx_in_no"`
	WhsIdf       int       `json:"whs_idf"`
	TrxInDate    string    `json:"trx_in_date"`
	TrxInSuppIdf int       `json:"trx_in_supp_idf"`
	TrxInNotes   string    `json:"trx_in_notes"`
	CreatedAt    time.Time `json:"created_at"`
}

type TransactionInDetail struct {
	TrxInDPK         int       `json:"trx_in_d_pk"`
	TrxInIdf         int       `json:"trx_in_idf"`
	TrxInDProductIdf int       `json:"trx_in_d_product_idf"`
	TrxInDQtyDus     int       `json:"trx_in_d_qty_dus"`
	TrxInDQtyPcs     int       `json:"trx_in_d_qty_pcs"`
	CreatedAt        time.Time `json:"created_at"`
}

type TransactionIn struct {
	Header  TransactionInHeader   `json:"header"`
	Details []TransactionInDetail `json:"details"`
}

type TransactionInList struct {
	Header        TransactionInHeader `json:"header"`
	WarehouseName string              `json:"warehouse_name"`
	SupplierName  string              `json:"supplier_name"`
}

type TransactionInDetailView struct {
	Detail      TransactionInDetail `json:"detail"`
	ProductName string              `json:"product_name"`
}

type CreateTransactionInRequest struct {
	TrxInNo      string                             `json:"trx_in_no"`
	WhsIdf       int                                `json:"whs_idf"`
	TrxInDate    string                             `json:"trx_in_date"`
	TrxInSuppIdf int                                `json:"trx_in_supp_idf"`
	TrxInNotes   string                             `json:"trx_in_notes"`
	Details      []CreateTransactionInDetailRequest `json:"details"`
}

type CreateTransactionInDetailRequest struct {
	TrxInDProductIdf int `json:"trx_in_d_product_idf"`
	TrxInDQtyDus     int `json:"trx_in_d_qty_dus"`
	TrxInDQtyPcs     int `json:"trx_in_d_qty_pcs"`
}
