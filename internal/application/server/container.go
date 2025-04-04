package server

import (
	"chatY-go/internal/api/tcp"
	"chatY-go/internal/application/server/session"
	"chatY-go/internal/domain/message"
	"chatY-go/internal/domain/user"
	"chatY-go/pkg/authclient"
	"chatY-go/pkg/config"
	"chatY-go/pkg/logger"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"sync"
)

const (
	SYS_USER = "System"
)

type IApplication interface {
	Init()
	AppConfig() config.IConfig
	Session() session.IChatServer
	ChatServer() tcp.IRunnable
	GRPCClient() authclient.IAuthClient
}

type Application struct {
	config     config.IConfig
	server     tcp.IRunnable
	session    session.IChatServer
	logger     logger.ILogger
	gRPCClient authclient.IAuthClient
}

func NewApplication(logger logger.ILogger) *Application {
	return &Application{logger: logger}
}

func (app *Application) AppConfig() config.IConfig {
	if app.config == nil {
		config := config.NewConfig()
		app.config = config
	}
	return app.config
}

func (app *Application) Logger() logger.ILogger {
	return app.logger
}

func (app *Application) Session() session.IChatServer {
	if app.session == nil {
		chatSession := session.NewChatServer(
			user.NewUser(SYS_USER),
			make(chan message.IMessage),
			make(map[net.Conn]user.IUser),
			&sync.Mutex{},
			app.Logger(),
			app.GRPCClient(),
		)
		app.session = chatSession

		app.logger.Info("[SERVER] Chat Session initialized")
	}
	return app.session
}

func (app *Application) GRPCClient() authclient.IAuthClient {
	if app.gRPCClient == nil {
		addr := fmt.Sprintf("%s:%s", app.AppConfig().ServerAddress(), app.AppConfig().AuthPort())
		authConn, _ := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		app.gRPCClient = authclient.NewAuthClient(authConn)
	}
	return app.gRPCClient
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
