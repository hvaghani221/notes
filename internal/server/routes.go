package server

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var signInKey = []byte("secret")

type jwtClaim struct {
	jwt.RegisteredClaims
	ID int32 `json:"id"`
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	users := e.Group("/api/auth")
	users.POST("/signup", s.CreateUser)
	users.POST("/login", s.LogIn)

	notes := e.Group("/api/notes")
	notes.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: signInKey,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtClaim)
		},
	}))
	notes.GET("/", s.ListNotes)
	notes.GET("/:id", s.GetNote)
	notes.POST("/", s.CreateNote)
	notes.PUT("/:id", s.UpdateNote)
	notes.DELETE("/:id", s.DeleteNote)
	notes.POST("/:id/share", s.ShareNote)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

// func (s *Server) healthHandler(c echo.Context) error {
// 	return c.JSON(http.StatusOK, s.db.Health())
// }
