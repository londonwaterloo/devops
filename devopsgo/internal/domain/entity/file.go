package entity

import "time"

// File представляет файл в системе
type File struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Content   []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
