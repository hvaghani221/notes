package model

import (
	"time"

	"notes/internal/database"
)

type User struct {
	ID        int32     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func UserFromDB(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt.Time,
	}
}
