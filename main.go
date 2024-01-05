package main

import (
	"fmt"

	"notes/internal/server"
)

func main() {
	config, err := server.LoadConfigFromEnv()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

	server := server.NewServer(config)

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
