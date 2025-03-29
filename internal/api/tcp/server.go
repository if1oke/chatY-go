package tcp

import (
	"chatY-go/internal/domain/session"
	"chatY-go/pkg/utils"
	"fmt"
	"net"
)

const (
	tcp = "tcp"
)

type IRunnable interface {
	Start(config utils.IConfig) error
}

type Server struct {
	session session.ISession
}

func NewServer(s session.ISession) *Server {
	return &Server{session: s}
}

func (s *Server) Start(config utils.IConfig) error {
	listen, err := net.Listen(tcp, fmt.Sprintf("%s:%s", config.ServerAddress(), config.ServerPort()))
	if err != nil {
		return fmt.Errorf("listen err: %s", err.Error())
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept err: %s", err.Error())
			continue
		}

		go s.session.Start(conn)
	}

}
