package main

import "chatY-go/internal/application/client"

func main() {
	//var log logger.ILogger
	//var app client_old.IApplication
	//
	//log = logger.NewClientLogger()
	//app = client_old.NewApplication(log)
	//if err := app.Init(); err != nil {
	//	sysLog.Fatalf("Error initializing application: %v", err)
	//}
	//
	//app.Session().Start()

	app := client.NewApplication()
	app.AskCredentials()

	err := client.Run(app)
	if err != nil {
		panic(err)
	}

}
