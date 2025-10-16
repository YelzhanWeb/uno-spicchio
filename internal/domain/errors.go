package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user is not active")
)

var ErrCategoryNotFound = errors.New("category not found")

var ErrDishNotFound = errors.New("dish not found")

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrInsufficientStock   = errors.New("insufficient stock for ingredient")
	ErrInvalidStatusChange = errors.New("invalid status change")
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)
