package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"notes/internal/database"
	"notes/internal/model"
)

type userCreateDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type logInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) CreateUser(c echo.Context) error {
	var userDTO userCreateDTO
	if err := c.Bind(&userDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// TODO: validation

	hash, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Printf("signup password hash: %s\n", string(hash))

	dbUser, err := s.db.CreateUser(c.Request().Context(), database.CreateUserParams{
		Username:     userDTO.Username,
		Email:        userDTO.Email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, model.UserFromDB(dbUser))
}

func (s *Server) LogIn(c echo.Context) error {
	var login logInDTO
	if err := c.Bind(&login); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fmt.Printf("login password hash: %s\n", string(hash))

	dbUser, err := s.db.GetUserByEmail(c.Request().Context(), login.Email)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
		// return echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordHash), []byte(login.Password)); err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
		// return echo.ErrUnauthorized
	}

	user := model.UserFromDB(dbUser)
	log.Println("user: ", user)

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
