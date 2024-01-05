package unittest

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"notes/internal/model"
)

type MockRepository struct {
	users       map[int32]model.User
	notes       map[int32]model.Note
	sharedNotes map[int32][]int32 // Map of noteID to a slice of userIDs who have access
	mu          sync.Mutex
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users:       make(map[int32]model.User),
		notes:       make(map[int32]model.Note),
		sharedNotes: make(map[int32][]int32),
	}
}

func (m *MockRepository) CreateUser(ctx context.Context, user model.UserCreateDTO) (*model.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	newUser := model.User{
		ID:           int32(len(m.users) + 1),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.Password, // Simulate password hash
		CreatedAt:    time.Now(),
	}

	m.users[newUser.ID] = newUser
	return &newUser, nil
}

func (m *MockRepository) GetUserByEmailAndPassword(ctx context.Context, login model.LogInDTO) (*model.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, user := range m.users {
		if user.Email == login.Email && user.PasswordHash == login.Password { // Simulate password hash check
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockRepository) CreateNote(ctx context.Context, noteDTO model.NoteDTO) (*model.Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	newNote := model.Note{
		ID:        int32(len(m.notes) + 1),
		UserID:    noteDTO.UserID,
		Title:     noteDTO.Title,
		Content:   noteDTO.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.notes[newNote.ID] = newNote
	return &newNote, nil
}

func (m *MockRepository) ListNotesByUserID(ctx context.Context, userID int32) ([]model.Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var notes []model.Note
	for _, note := range m.notes {
		if note.UserID == userID {
			notes = append(notes, note)
		}
	}
	return notes, nil
}

func (m *MockRepository) GetNoteByUserID(ctx context.Context, noteID, userID int32) (*model.Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	note, ok := m.notes[noteID]
	if !ok || (note.UserID != userID && !m.isNoteSharedWithUser(noteID, userID)) {
		return nil, errors.New("note not found or access denied")
	}
	return &note, nil
}

func (m *MockRepository) UpdateNote(ctx context.Context, noteID int32, noteDTO model.NoteDTO) (*model.Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	note, ok := m.notes[noteID]
	if !ok || (note.UserID != noteDTO.UserID && !m.isNoteSharedWithUser(noteID, noteDTO.UserID)) {
		return nil, errors.New("note not found or access denied")
	}

	updatedNote := model.Note{
		ID:        noteID,
		UserID:    noteDTO.UserID,
		Title:     noteDTO.Title,
		Content:   noteDTO.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: time.Now(),
	}

	m.notes[noteID] = updatedNote
	return &updatedNote, nil
}

func (m *MockRepository) DeleteNote(ctx context.Context, noteID, userID int32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	note, ok := m.notes[noteID]
	if !ok || (note.UserID != userID && !m.isNoteSharedWithUser(noteID, userID)) {
		return errors.New("note not found or access denied")
	}

	delete(m.notes, noteID)
	return nil
}

func (m *MockRepository) ShareNote(ctx context.Context, noteShareDTO model.NoteShareDTO) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, noteExists := m.notes[noteShareDTO.NoteID]
	if !noteExists {
		return errors.New("note not found")
	}

	m.sharedNotes[noteShareDTO.NoteID] = append(m.sharedNotes[noteShareDTO.NoteID], noteShareDTO.UserID)
	return nil
}

func (m *MockRepository) SearchNotes(ctx context.Context, userID int32, query string) ([]model.Note, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var notes []model.Note
	for _, note := range m.notes {
		if note.UserID == userID && contains(note.Content, query) {
			notes = append(notes, note)
		}
	}
	return notes, nil
}

// contains checks if the text contains the query (case-insensitive)
func contains(text, query string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(query))
}

func (m *MockRepository) isNoteSharedWithUser(noteID, userID int32) bool {
	for _, id := range m.sharedNotes[noteID] {
		if id == userID {
			return true
		}
	}
	return false
}
