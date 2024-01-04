package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"notes/internal/model"
)

func (s *Server) CreateUser(c echo.Context) error {
	var userDTO model.UserCreateDTO
	if err := c.Bind(&userDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// TODO: validation

	hash, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Printf("signup password hash: %s\n", string(hash))

	user, err := s.repository.CreateUser(c.Request().Context(), userDTO)
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

	user, err := s.repository.GetUserByEmailAndPassword(c.Request().Context(), login)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	claims := &jwtClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		ID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signInKey)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}
