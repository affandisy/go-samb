package model

import "time"

type Warehouse struct {
	WhsPK     int       `json:"whs_pk"`
	WhsName   string    `json:"whs_name"`
	CreatedAt time.Time `json:"created_at"`
}
