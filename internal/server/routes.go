package server

import (
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

type jwtClaim struct {
	jwt.RegisteredClaims
	ID int32 `json:"id"`
}

func (s *Server) RegisterRoutes() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(s.Config.RateLimit))))

	users := e.Group("/api/auth")
	users.POST("/signup", s.CreateUser)
	users.POST("/login", s.LogIn)

	notes := e.Group("/api/notes")
	notes.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(s.Config.SignInKey),
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

	notes.GET("/search", s.SearchNotes)

	return e
}
