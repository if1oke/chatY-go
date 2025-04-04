package main

import (
	"chatY-go/internal/application/client"
	"chatY-go/pkg/logger"
	sysLog "log"
)

func main() {
	var log logger.ILogger
	var app client.IApplication

	log = logger.NewClientLogger()
	app = client.NewApplication(log)
	if err := app.Init(); err != nil {
		sysLog.Fatalf("Error initializing application: %v", err)
	}

	app.Session().Start()
}
