package domain

type Ingredient struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Unit   string  `json:"unit"`
	Qty    float64 `json:"qty"`
	MinQty float64 `json:"min_qty"`
}

func (i *Ingredient) IsLowStock() bool {
	return i.Qty <= i.MinQty
}
