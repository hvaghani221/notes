package model

import (
	"time"

	"notes/internal/database"
)

type Note struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NoteFromDB(dbNote database.Note) Note {
	return Note{
		ID:        dbNote.ID,
		UserID:    dbNote.UserID.Int32,
		Title:     dbNote.Title,
		Content:   dbNote.Content,
		CreatedAt: dbNote.CreatedAt.Time,
		UpdatedAt: dbNote.UpdatedAt.Time,
	}
}
