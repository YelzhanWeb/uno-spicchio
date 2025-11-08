package domain

import "errors"

// Auth errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user is not active")
)

// Category errors
var ErrCategoryNotFound = errors.New("category not found")

// Dish errors
var ErrDishNotFound = errors.New("dish not found")

// Order errors
var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrInsufficientStock   = errors.New("insufficient stock for ingredient")
	ErrInvalidStatusChange = errors.New("invalid status change")
)

// User errors
var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

// Table errors
var (
	ErrTableNotFound = errors.New("table not found")
	ErrTableBusy     = errors.New("table is busy")
)

// Ingredient errors
var ErrIngredientNotFound = errors.New("ingredient not found")

// Supply errors
var ErrSupplyNotFound = errors.New("supply not found")
