package client

import (
	"bufio"
	"chatY-go/internal/application/client/session"
	"chatY-go/pkg/config"
	"fmt"
	"os"
	"strings"
)

type IClientConfig interface {
	Username() string
	SetUsername(username string)
	Password() string
	SetPassword(password string)
	ServerAddress() string
	ServerPort() string
}

type IClientSession interface {
	Authenticate() error
	Reader() *bufio.Reader
	Writer() *bufio.Writer
}

type Application struct {
	config  IClientConfig
	session IClientSession
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) Config() IClientConfig {
	if a.config == nil {
		a.config = config.NewClientConfig()
	}
	return a.config
}

func (a *Application) Session() IClientSession {
	if a.session == nil {
		addr := fmt.Sprintf("%s:%s", a.config.ServerAddress(), a.config.ServerPort())
		s, err := session.NewSession(addr, a.Config().Username(), a.Config().Password())
		if err != nil {
			panic(err)
		}
		a.session = s
	}
	return a.session
}

func (a *Application) AskCredentials() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter your username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Please enter your password: ")
	password, _ := reader.ReadString('\n')

	a.Config().SetUsername(strings.TrimRight(username, "\n"))
	a.Config().SetPassword(strings.TrimRight(password, "\n"))
}
