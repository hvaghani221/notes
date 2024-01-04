package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"

	"notes/internal/database/generated"
	"notes/internal/model"
	"notes/internal/repository"
)

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

type Repository struct {
	db      *sql.DB
	queries *generated.Queries
}

var _ repository.Repository = (*Repository)(nil)

func (r *Repository) CreateUser(ctx context.Context, user model.UserCreateDTO) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	dbUser, err := r.queries.CreateUser(ctx, generated.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return model.User{}, err
	}

	return dbUserToUser(dbUser), nil
}

func dbUserToUser(dbUser generated.User) model.User {
	return model.User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		CreatedAt:    dbUser.CreatedAt.Time,
	}
}

func dbNoteToNote(dbNote generated.Note) model.Note {
	return model.Note{
		ID:        dbNote.ID,
		UserID:    dbNote.UserID,
		Title:     dbNote.Title,
		Content:   dbNote.Content,
		CreatedAt: dbNote.CreatedAt,
		UpdatedAt: dbNote.UpdatedAt,
	}
}

func (r *Repository) GetUserByEmailAndPassword(ctx context.Context, login model.LogInDTO) (model.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return model.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(login.Password)); err != nil {
		return model.User{}, err
	}
	return dbUserToUser(dbUser), nil
}

func (r *Repository) CreateNote(ctx context.Context, note model.NoteDTO) (model.Note, error) {
	dbNote, err := r.queries.CreatNote(ctx, generated.CreatNoteParams{
		UserID:  note.UserID,
		Title:   note.Title,
		Content: note.Content,
	})
	if err != nil {
		return model.Note{}, err
	}

	return dbNoteToNote(dbNote), nil
}

func (r *Repository) ListNotesByUserID(ctx context.Context, userID int32) ([]model.Note, error) {
	dbNotes, err := r.queries.ListNotesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	notes := make([]model.Note, 0, len(dbNotes))

	for _, dbNote := range dbNotes {
		notes = append(notes, dbNoteToNote(dbNote))
	}

	return notes, nil
}

func (r *Repository) GetNoteByUserID(ctx context.Context, noteID int32, userID int32) (model.Note, error) {
	dbNote, err := r.queries.GetNoteByUserID(ctx, generated.GetNoteByUserIDParams{
		ID:     noteID,
		UserID: userID,
	})
	if err != nil {
		return model.Note{}, err
	}

	return dbNoteToNote(dbNote), nil
}

func (r *Repository) UpdateNote(ctx context.Context, noteID int32, note model.NoteDTO) (model.Note, error) {
	fmt.Println("updating note", noteID, note)
	dbNote, err := r.queries.UpdateNote(ctx, generated.UpdateNoteParams{
		ID:      noteID,
		UserID:  note.UserID,
		Title:   note.Title,
		Content: note.Content,
	})
	if err != nil {
		return model.Note{}, err
	}
	return dbNoteToNote(dbNote), nil
}

func (r *Repository) DeleteNote(ctx context.Context, noteID, userID int32) error {
	return r.queries.DeleteNote(ctx, generated.DeleteNoteParams{
		ID:     noteID,
		UserID: userID,
	})
}

func (r *Repository) ShareNote(ctx context.Context, share model.NoteShareDTO) error {
	fmt.Println("sharing note", share)
	_, err := r.queries.ShareNote(ctx, generated.ShareNoteParams{
		Noteid:          share.NoteID,
		Userid:          share.UserID,
		Sharedwithemail: share.SharedWith,
	})

	return err
}

func (r *Repository) SearchNotes(ctx context.Context, id int32, query string) ([]model.Note, error) {
	fmt.Println("searching notes", id, query)
	dbNotes, err := r.queries.SearchNotes(ctx, generated.SearchNotesParams{
		UserID:         id,
		PlaintoTsquery: query,
	})
	if err != nil {
		return nil, err
	}

	notes := make([]model.Note, 0, len(dbNotes))

	for _, dbNote := range dbNotes {
		notes = append(notes, dbNoteToNote(dbNote))
	}

	return notes, nil
}

func New() *Repository {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	return &Repository{
		db:      db,
		queries: generated.New(db),
	}
}
