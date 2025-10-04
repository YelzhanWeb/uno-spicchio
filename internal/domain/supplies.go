package domain

import "time"

type Supply struct {
	ID           int       `json:"id"`
	IngredientID int       `json:"ingredient_id"`
	Qty          float64   `json:"qty"`
	SupplierName string    `json:"supplier_name"`
	CreatedAt    time.Time `json:"created_at"`
}
