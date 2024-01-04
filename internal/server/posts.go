package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"notes/internal/database"
	"notes/internal/model"
)

type noteDTO struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (s *Server) ListNotes(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID

	notes, err := s.db.ListNotesByUserID(c.Request().Context(), userID)
	if err != nil {
		return echo.ErrInternalServerError
	}

	if len(notes) == 0 {
		return c.JSON(http.StatusOK, []model.Note{})
	}

	return c.JSON(http.StatusOK, notes)
}

func (s *Server) GetNote(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	dbNote, err := s.db.GetNoteByUserID(c.Request().Context(), database.GetNoteByUserIDParams{
		ID:     int32(noteID),
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, model.NoteFromDB(dbNote))
}

func (s *Server) CreateNote(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID

	var noteDTO noteDTO
	if err := c.Bind(&noteDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	dbNote, err := s.db.CreatNote(c.Request().Context(), database.CreatNoteParams{
		UserID:  userID,
		Title:   noteDTO.Title,
		Content: noteDTO.Content,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.NoteFromDB(dbNote))
}

func (s *Server) UpdateNote(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	var noteDTO noteDTO
	if err := c.Bind(&noteDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	dbNote, err := s.db.UpdateNote(c.Request().Context(), database.UpdateNoteParams{
		ID:      int32(noteID),
		UserID:  userID,
		Title:   noteDTO.Title,
		Content: noteDTO.Content,
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.NoteFromDB(dbNote))
}

func (s *Server) DeleteNote(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	err = s.db.DeleteNote(c.Request().Context(), database.DeleteNoteParams{
		ID:     int32(noteID),
		UserID: userID,
	})

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
