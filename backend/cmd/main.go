package main

import (
	"Voice-Assistant/internal/app"
	"Voice-Assistant/internal/db"
	"Voice-Assistant/internal/handler"
	"log"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	handler.InitHandlers(database)
	server := app.NewServer()
	server.RunWithGracefulShutdown()
}
