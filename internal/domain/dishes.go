package domain

type Dish struct {
	ID          int       `json:"id"`
	CategoryID  int       `json:"category_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Price       float64   `json:"price"`
	PhotoURL    *string   `json:"photo_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	Category    *Category `json:"category,omitempty"`
}

type DishIngredient struct {
	DishID       int         `json:"dish_id"`
	IngredientID int         `json:"ingredient_id"`
	QtyPerDish   float64     `json:"qty_per_dish"`
	Ingredient   *Ingredient `json:"ingredient,omitempty"`
}
