package server

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"notes/internal/database"
	"notes/internal/model"
)

type createNoteDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *Server) ListNotes(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID

	log.Println("userID: ", userID)

	notes, err := s.db.GetNotesByUserID(c.Request().Context(), pgtype.Int4{Int32: userID, Valid: true})
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(notes) == 0 {
		return c.JSON(http.StatusOK, []model.Note{})
	}

	return c.JSON(http.StatusOK, notes)
}

func (s *Server) GetNote(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) CreateNote(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID

	var noteDTO createNoteDTO
	if err := c.Bind(&noteDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	dbNote, err := s.db.CreatNote(c.Request().Context(), database.CreatNoteParams{
		UserID:  pgtype.Int4{Int32: userID, Valid: true},
		Title:   noteDTO.Title,
		Content: noteDTO.Content,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.NoteFromDB(dbNote))
}

func (s *Server) UpdateNote(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) DeleteNote(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}
