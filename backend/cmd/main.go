package main

import (
	"log"
	"os"

	"github.com/user/taskapp/internal/app"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			app.RunMigrations()
			return
		case "seed":
			app.SeedDatabase()
			return
		}
	}

	server := app.NewServer()
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
