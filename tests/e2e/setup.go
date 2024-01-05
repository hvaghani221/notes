package e2e

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gavv/httpexpect/v2"
	"github.com/joho/godotenv"

	"notes/internal/database"
	"notes/internal/model"
)

var dbname, password, username, port, host, url string

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbname = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port = os.Getenv("DB_PORT")
	host = os.Getenv("DB_HOST")
	url = fmt.Sprintf("http://%s:%s", "localhost", os.Getenv("PORT"))
}

var users = []model.UserCreateDTO{
	{
		Username: "username1",
		Email:    "user1@email.com",
		Password: "Hello@123",
	},
	{
		Username: "username2",
		Email:    "user2@email.com",
		Password: "Hello@123",
	},
	{
		Username: "username3",
		Email:    "user3@email.com",
		Password: "Hello@123",
	},
}

func setup(e *httpexpect.Expect) error {
	repo := database.NewFromConfig(&database.Config{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: dbname,
	})
	defer repo.Db.Close()

	if err := cleanup(repo.Db); err != nil {
		return fmt.Errorf("error cleaning up: %v", err)
	}

	for _, user := range users {
		e.POST("/api/auth/signup").WithJSON(user).Expect().Status(http.StatusOK)
		// if _, err := repo.CreateUser(context.Background(), user); err != nil {
		// 	return fmt.Errorf("error creating user: %v", err)
		// }
	}

	return nil
}

func cleanup(db *sql.DB) error {
	tables := []string{"users", "notes", "shared_notes"}
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;")
		if err != nil {
			return fmt.Errorf("error truncating table %s: %v", table, err)
		}
	}

	return nil
}
