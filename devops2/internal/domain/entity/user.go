package entity

import "time"

// User представляет пользователя в системе
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"` // пароль не возвращается в JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
