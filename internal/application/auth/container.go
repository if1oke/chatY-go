package auth

import (
	"chatY-go/internal/api/grpc/auth"
	"chatY-go/internal/api/grpc/auth/proto"
	auth2 "chatY-go/internal/domain/auth"
	"chatY-go/pkg/config"
	"chatY-go/pkg/logger"
	"google.golang.org/grpc"
)

type IApplication interface {
	AuthHandler() proto.AuthServiceServer
	GRPCServer() *grpc.Server
	Config() config.IConfig
	Init()
}

type Application struct {
	config     config.IConfig
	logger     logger.ILogger
	gRPCServer *grpc.Server
	service    auth2.IAuthService
	handler    proto.AuthServiceServer
}

func NewApplication(l logger.ILogger) *Application {
	return &Application{logger: l}
}

func (a *Application) Config() config.IConfig {
	if a.config == nil {
		a.config = config.GetConfig()
	}
	return a.config
}

func (a *Application) GRPCServer() *grpc.Server {
	if a.gRPCServer == nil {
		a.gRPCServer = grpc.NewServer()
	}
	return a.gRPCServer
}

func (a *Application) AuthService() auth2.IAuthService {
	if a.service == nil {
		a.service = NewAuthService(a.Logger())
	}
	return a.service
}

func (a *Application) Logger() logger.ILogger {
	return a.logger
}

func (a *Application) AuthHandler() proto.AuthServiceServer {
	if a.handler == nil {
		a.handler = auth.NewAuthHandler(a.AuthService())
	}
	return a.handler
}

func (a *Application) Init() {
	a.config = a.Config()
	a.gRPCServer = a.GRPCServer()
	a.service = a.AuthService()
	a.handler = a.AuthHandler()
}
