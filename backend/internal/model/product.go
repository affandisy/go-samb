package model

import "time"

type Product struct {
	ProductPK   int       `json:"product_pk"`
	ProductName string    `json:"product_name"`
	CreatedAt   time.Time `json:"created_at"`
}
