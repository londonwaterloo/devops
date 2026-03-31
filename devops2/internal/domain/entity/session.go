package entity

import "time"

// Session представляет сессию пользователя
type Session struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
