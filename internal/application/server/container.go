package server

import (
	"chatY-go/internal/api/tcp"
	appSession "chatY-go/internal/application/session"
	"chatY-go/internal/domain/message"
	"chatY-go/internal/domain/session"
	"chatY-go/internal/domain/user"
	"chatY-go/pkg/config"
	"chatY-go/pkg/logger"
	"net"
	"sync"
)

const (
	SYS_USER = "System"
)

type IApplication interface {
	Init()
	AppConfig() config.IConfig
	Session() session.ISession
	ChatServer() tcp.IRunnable
}

type Application struct {
	config  config.IConfig
	server  tcp.IRunnable
	session session.ISession
	logger  logger.ILogger
}

func NewApplication(logger logger.ILogger) *Application {
	return &Application{logger: logger}
}

func (app *Application) AppConfig() config.IConfig {
	if app.config == nil {
		config := config.GetConfig()
		app.config = config
	}
	return app.config
}

func (app *Application) Logger() logger.ILogger {
	return app.logger
}

func (app *Application) Session() session.ISession {
	if app.session == nil {
		chatSession := appSession.NewChatSession(
			user.NewUser(SYS_USER),
			make(chan message.IMessage),
			make(map[net.Conn]user.IUser),
			&sync.Mutex{},
			app.Logger(),
		)
		app.session = chatSession

		app.logger.Info("[SERVER] Chat Session initialized")
	}
	return app.session
}

func (app *Application) ChatServer() tcp.IRunnable {
	if app.server == nil {
		chatServer := tcp.NewServer(app.Session(), app.Logger())
		app.server = chatServer

		app.logger.Info("[SERVER] TCP Server initialized")
	}
	return app.server
}

func (app *Application) Init() {
	app.config = app.AppConfig()
	app.server = app.ChatServer()
}
