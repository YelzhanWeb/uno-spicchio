package domain

type Dish struct {
	ID          int     `json:"id" db:"id"`
	CategoryID  int     `json:"category_id" db:"category_id"` // <-- ИСПРАВЛЕНО
	Name        string  `json:"name" db:"name"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	PhotoURL    string  `json:"photo_url" db:"photo_url"`
	IsActive    bool    `json:"is_active" db:"is_active"`
}

type DishIngredient struct {
	DishID       int     `json:"dish_id"`
	IngredientID int     `json:"ingredient_id"`
	QtyPerDish   float64 `json:"qty_per_dish"`
}
