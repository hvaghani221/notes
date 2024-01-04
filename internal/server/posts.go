package server

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"notes/internal/model"
)

func (s *Server) ListNotes(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID

	notes, err := s.repository.ListNotesByUserID(c.Request().Context(), userID)
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

	note, err := s.repository.GetNoteByUserID(c.Request().Context(), int32(noteID), userID)
	if err != nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, note)
}

func (s *Server) CreateNote(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID

	var noteDTO model.NoteDTO
	if err := c.Bind(&noteDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	noteDTO.UserID = userID

	note, err := s.repository.CreateNote(c.Request().Context(), noteDTO)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, note)
}

func (s *Server) UpdateNote(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	var noteDTO model.NoteDTO
	if err := c.Bind(&noteDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	noteDTO.UserID = userID

	note, err := s.repository.UpdateNote(c.Request().Context(), int32(noteID), noteDTO)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, note)
}

func (s *Server) DeleteNote(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	err = s.repository.DeleteNote(c.Request().Context(), int32(noteID), userID)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
