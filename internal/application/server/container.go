package server

import (
	"chatY-go/internal/api/tcp"
	appSession "chatY-go/internal/application/session"
	"chatY-go/internal/domain/message"
	"chatY-go/internal/domain/session"
	"chatY-go/internal/domain/user"
	"chatY-go/pkg/utils"
	"net"
	"sync"
)

const (
	SYS_USER = "System"
)

type IApplication interface {
	Init()
	AppConfig() utils.IConfig
	Session() session.ISession
	ChatServer() tcp.IRunnable
}

type Application struct {
	config  utils.IConfig
	server  tcp.IRunnable
	session session.ISession
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) AppConfig() utils.IConfig {
	if app.config == nil {
		config := utils.GetConfig()
		app.config = config
	}
	return app.config
}

func (app *Application) Session() session.ISession {
	if app.session == nil {
		chatSession := appSession.NewChatSession(
			user.NewUser(SYS_USER),
			make(chan message.IMessage),
			make(map[net.Conn]user.IUser),
			&sync.Mutex{},
		)
		app.session = chatSession
	}
	return app.session
}

func (app *Application) ChatServer() tcp.IRunnable {
	if app.server == nil {
		chatServer := tcp.NewServer(app.Session())
		app.server = chatServer
	}
	return app.server
}

func (app *Application) Init() {
	app.config = app.AppConfig()
	app.server = app.ChatServer()
}
