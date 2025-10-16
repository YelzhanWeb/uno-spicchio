// file: internal/domain/category.go

package domain

// Category представляет категорию блюд.
// type Category struct {
// 	ID     int    `json:"id" db:"id"`
// 	Name   string `json:"name" db:"name"`
// 	Dishes []Dish `json:"dishes,omitempty"` // Поле для хранения списка блюд
// }

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
