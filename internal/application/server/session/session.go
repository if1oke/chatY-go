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
	}()

	// Ask credentials
	login, err := s.AskUser("Input username:", conn)
	if err != nil {
		s.logger.Errorf("[DISCONNECT] Error ask login: %v", err)
	}

	password, err := s.AskUser("Input password:", conn)
	if err != nil {
		s.logger.Errorf("[DISCONNECT] Error ask password: %v", err)
	}

	ok, msg, err := s.authClient.Login(login, password)
	if err != nil || !ok {
		s.logger.Errorf("[DISCONNECT] Error login: %s, %v", msg, err)
		conn.Close()
		return
	}
	s.logger.Infof("[AUTH] Пользователь %s успешно авторизован", login)

	s.register(conn, login)
	s.logger.Infof("[JOIN] %s joined the chat", s.clients[conn].Nickname())

	defer func() {
		s.unregister(conn)
		err := conn.Close()
		if err != nil {
			return
		}
	}()

	reader := bufio.NewReader(conn)

	for {
		rawMessage, err := reader.ReadString('\n')

		if err != nil {
			s.unregister(conn)
			return
		}

		msg := message.NewMessage(s.clients[conn], rawMessage)
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
				s.logger.Errorf("[ERROR] Write to client %s failed: %v", client.RemoteAddr(), err)
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
