package client_old

import (
	"chatY-go/pkg/config"
	"chatY-go/pkg/logger"
)

type IApplication interface {
	Init() error
	Session() IClientSession
	Logger() logger.ILogger
}

type Application struct {
	config  config.IClientConfig
	logger  logger.ILogger
	session IClientSession
}

func NewApplication(logger logger.ILogger) *Application {
	return &Application{
		logger: logger,
	}
}

func (a *Application) Config() config.IClientConfig {
	if a.config == nil {
		a.config = config.NewClientConfig()
	}
	return a.config
}

func (a *Application) Logger() logger.ILogger {
	return a.logger
}

func (a *Application) Session() IClientSession {
	if a.session == nil {
		a.session = NewClientSession(a.Logger(), a.Config())
	}
	return a.session
}

func (a *Application) Init() error {
	a.config = a.Config()
	a.session = a.Session()
	return nil
}
