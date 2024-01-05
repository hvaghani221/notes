package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"notes/internal/model"
	"notes/internal/validator"
)

func (s *Server) CreateUser(c echo.Context) error {
	var userDTO model.UserCreateDTO
	if err := c.Bind(&userDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := validator.Email(userDTO.Email); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := validator.Password(userDTO.Password); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user, err := s.Repository.CreateUser(c.Request().Context(), userDTO)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (s *Server) LogIn(c echo.Context) error {
	var login model.LogInDTO
	if err := c.Bind(&login); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := validator.Email(login.Email); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	user, err := s.Repository.GetUserByEmailAndPassword(c.Request().Context(), login)
	if err != nil {
		c.Logger().Error(fmt.Errorf("failed to get user[%s]: %w", login.Email, err))
		return echo.ErrUnauthorized
	}

	claims := &jwtClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		ID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.Config.SignInKey))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}
