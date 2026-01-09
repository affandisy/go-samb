package model

import "time"

type Customer struct {
	CustomerPK   int       `json:"customer_pk"`
	CustomerName string    `json:"customer_name"`
	CreatedAt    time.Time `json:"created_at"`
}
