package main

import (
	"chatY-go/internal/application/server"
	"chatY-go/pkg/logger"
	"log"
)

func main() {
	var app server.IApplication
	var logger = logger.NewLogger()

	app = server.NewApplication(logger)
	app.Init()

	err := app.ChatServer().Start(app.AppConfig())
	if err != nil {
		log.Fatalf("server start err: %v", err)
	}

}
