package session

import (
	"bufio"
	"chatY-go/internal/domain/message"
	"chatY-go/internal/domain/user"
	"chatY-go/pkg/authclient"
	"chatY-go/pkg/logger"
	"fmt"
	"net"
	"strings"
	"sync"
)

const (
	CommandNickname = "/nick"
	CommandList     = "/list"
	CommandExit     = "/exit"
	CommandHelp     = "/help"
	CommandWhisper  = "/whisper"
	CommandAuth     = "/auth"
)

type IChatServer interface {
	Start(conn net.Conn)
}

type ChatServer struct {
	systemUser user.IUser
	broadcast  chan message.IMessage
	clients    map[net.Conn]user.IUser
	mu         *sync.Mutex
	logger     logger.ILogger
	authClient authclient.IAuthClient
}

func NewChatServer(
	systemUser user.IUser,
	broadcast chan message.IMessage,
	clients map[net.Conn]user.IUser,
	mu *sync.Mutex,
	logger logger.ILogger,
	authClient authclient.IAuthClient,
) *ChatServer {
	s := &ChatServer{
		systemUser: systemUser,
		broadcast:  broadcast,
		clients:    clients,
		mu:         mu,
		logger:     logger,
		authClient: authClient,
	}
	go s.broadcaster()
	return s
}

func (s *ChatServer) AskUser(msg string, conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	_, err := fmt.Fprintf(conn, msg)
	if err != nil {
		s.logger.Errorf("[DISCONNECT] Error ask %s: %v", msg, err)
		return "", fmt.Errorf("[DISCONNECT] Error ask %s: %v", msg, err)
	}

	prepVal, _ := reader.ReadString('\n')
	fmt.Fprint(conn, "\n")
	return strings.TrimSpace(prepVal), nil
}

func (s *ChatServer) Start(conn net.Conn) {
	defer func() {
		s.logger.Infof("[DISCONNECT] Client %s disconnected", conn.RemoteAddr())
		s.unregister(conn)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	var username, password string
	isAuthorized := false

	for {
		raw, err := reader.ReadString('\n')
		if err != nil {
			s.logger.Errorf("[ERROR] Error reading from %s: %v", conn.RemoteAddr(), err)
			return
		}
		cmd, args := parseCommand(raw)
		if !isAuthorized {
			if cmd != CommandAuth {
				s.sendMessageToConn(conn, "[ERROR] Please authorize using /auth <username> <password>\n")
				continue
			}

			if len(args) != 2 {
				s.sendMessageToConn(conn, "[ERROR] You must provide username and password\n")
			}

			username = args[0]
			password = args[1]

			ok, msg, err := s.authClient.Login(username, password)
			if err != nil || !ok {
				s.sendMessageToConn(conn, fmt.Sprintf("[AUTH] Failed to login: %v\n", msg))
			}

			s.register(conn, username)
			isAuthorized = true

			s.sendMessageToConn(conn, fmt.Sprintf("[AUTH] Successfully logged in: %v\n", username))
			s.logger.Infof("[JOIN] User %s joined the chat", username)

			continue
		}

		msg := message.NewMessage(s.clients[conn], raw)
		s.logger.Infof("[MESSAGE] %s", msg.Print())
		s.handleMessage(msg, conn)
	}
}

func (s *ChatServer) broadcaster() {
	for {
		msg := <-s.broadcast
		for client := range s.clients {
			_, err := fmt.Fprintf(client, msg.Print())
			if err != nil {
				s.logger.Errorf("[ERROR] Write to client_old %s failed: %v", client.RemoteAddr(), err)
			}
		}
	}
}

func (s *ChatServer) register(conn net.Conn, login string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[conn] = user.NewUser(login)
	s.notify(fmt.Sprintf("User %s joined the chat\n", s.clients[conn].Nickname()))
}

func (s *ChatServer) unregister(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, conn)
	conn.Close()
}
