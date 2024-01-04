package repository

import (
	"context"

	"notes/internal/model"
)

type UserRepository interface {
	CreateUser(context.Context, model.UserCreateDTO) (model.User, error)
	GetUserByEmailAndPassword(context.Context, model.LogInDTO) (model.User, error)
}

type NoteRepository interface {
	CreateNote(context.Context, model.NoteDTO) (model.Note, error)
	ListNotesByUserID(context.Context, int32) ([]model.Note, error)
	GetNoteByUserID(ctx context.Context, noteID, userID int32) (model.Note, error)
	UpdateNote(context.Context, int32, model.NoteDTO) (model.Note, error)
	DeleteNote(ctx context.Context, noteID, userID int32) error
}

type Repository interface {
	UserRepository
	NoteRepository
}
