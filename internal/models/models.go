package models

import "time"

type User struct {
	ID        int
	Username  string
	Password  string
	Role      string
	CreatedAt time.Time
}
