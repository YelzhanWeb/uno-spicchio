package domain

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleWaiter  Role = "waiter"
	RoleCook    Role = "cook"
	RoleManager Role = "manager"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	PhotoKey     string    `json:"photo_key"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}
