package model

import (
	"time"
)

type User struct {
	ID           int32     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	PasswordHash string    `json:"-"`
}

type UserCreateDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
