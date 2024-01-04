package model

import (
	"time"
)

type Note struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteDTO struct {
	UserID  int32  `json:"-"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
