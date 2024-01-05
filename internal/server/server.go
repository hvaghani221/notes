package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"notes/internal/database"
	"notes/internal/repository"
)

type Config struct {
	Host      string
	Port      int
	RateLimit int
	SignInKey string
}

func NewConfig(host string, port int, rateLimit int, signInKey string) Config {
	return Config{
		Host:      host,
		Port:      port,
		RateLimit: rateLimit,
		SignInKey: signInKey,
	}
}

func LoadConfigFromEnv() (Config, error) {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse port: %w", err)
	}
	rateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse port: %w", err)
	}
	signInKey := os.Getenv("SIGNIN_KEY")
	if signInKey == "" {
		return Config{}, fmt.Errorf("missing SIGNIN_KEY env")
	}
	if signInKey == "" {
		return Config{}, fmt.Errorf("missing SIGNIN_KEY env")
	}

	return Config{
		Host:      host,
		Port:      port,
		RateLimit: rateLimit,
		SignInKey: signInKey,
	}, nil
}

type Server struct {
	Config     Config
	Repository repository.Repository
}

func NewServer(config Config) *http.Server {
	NewServer := &Server{
		Config:     config,
		Repository: database.New(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.Config.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
