package model

type StockReport struct {
	Warehouse string `json:"warehouse"`
	Product   string `json:"product"`
	QtyDus    int    `json:"qty_dus"`
	QtyPcs    int    `json:"qty_pcs"`
}
