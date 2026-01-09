package model

import "time"

type Supplier struct {
	SupplierPK   int       `json:"supplier_pk"`
	SupplierName string    `json:"supplier_name"`
	CreatedAt    time.Time `json:"created_at"`
}
