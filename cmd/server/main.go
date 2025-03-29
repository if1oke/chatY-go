package main

import (
	"chatY-go/internal/application/server"
	"log"
)

func main() {
	var app server.IApplication

	app = server.NewApplication()
	app.Init()

	err := app.ChatServer().Start(app.AppConfig())
	log.Printf("Server started")
	if err != nil {
		log.Fatalf("server start err: %v", err)
	}

}
