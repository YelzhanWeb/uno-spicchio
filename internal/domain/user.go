package domain

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // скрываем при JSON
	Role      string    `json:"role"`
	PhotoKey  string    `json:"photo_key"` // ключ в MinIO, а не URL
	PhotoURL  string    `json:"photo_url"` // временный presigned URL (генерируется на лету)
	CreatedAt time.Time `json:"created_at"`
}
