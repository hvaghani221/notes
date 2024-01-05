package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"notes/internal/model"
	"notes/internal/repository"
	"notes/internal/validator"
)

func (s *Server) ListNotes(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID

	notes, err := s.Repository.ListNotesByUserID(c.Request().Context(), userID)
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

	note, err := s.Repository.GetNoteByUserID(c.Request().Context(), int32(noteID), userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.ErrNotFound
		}
		c.Logger().Error(fmt.Errorf("failed to get note[%d]: %w", noteID, err))
		return echo.ErrInternalServerError
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

	note, err := s.Repository.CreateNote(c.Request().Context(), noteDTO)
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

	note, err := s.Repository.UpdateNote(c.Request().Context(), int32(noteID), noteDTO)
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

	err = s.Repository.DeleteNote(c.Request().Context(), int32(noteID), userID)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.ErrForbidden
		}
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) ShareNote(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID
	noteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}

	var noteShareDTO model.NoteShareDTO
	if err := c.Bind(&noteShareDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := validator.Email(noteShareDTO.SharedWith); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	noteShareDTO.NoteID = int32(noteID)
	noteShareDTO.UserID = userID

	err = s.Repository.ShareNote(c.Request().Context(), noteShareDTO)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (s *Server) SearchNotes(c echo.Context) error {
	userID := c.Get("user").(*jwt.Token).Claims.(*jwtClaim).ID
	query := c.QueryParam("q")

	notes, err := s.Repository.SearchNotes(c.Request().Context(), userID, query)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, notes)
}
