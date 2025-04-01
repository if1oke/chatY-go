package tcp

import (
	"chatY-go/internal/application/server/session"
	"chatY-go/pkg/config"
	"chatY-go/pkg/logger"
	"fmt"
	"net"
)

const (
	tcp = "tcp"
)

type IRunnable interface {
	Start(config config.IConfig) error
}

type Server struct {
	session session.IChatServer
	logger  logger.ILogger
}

func NewServer(s session.IChatServer, l logger.ILogger) *Server {
	return &Server{session: s, logger: l}
}

func (s *Server) Start(config config.IConfig) error {
	listen, err := net.Listen(tcp, fmt.Sprintf("%s:%s", config.ServerAddress(), config.ServerPort()))
	if err != nil {
		s.logger.Errorf("[SERVER] Listen error: %v", err)
		return fmt.Errorf("listen err: %s", err.Error())
	}
	defer listen.Close()

	s.logger.Infof("[SERVER] Listening on %s:%s", config.ServerAddress(), config.ServerPort())

	for {
		conn, err := listen.Accept()
		if err != nil {
			s.logger.Warnf("[SERVER] Failed to accept connection: %v", err)
			continue
		}
		s.logger.Infof("[CONNECT] New connection from %s", conn.RemoteAddr())

		go s.session.Start(conn)
	}

}
