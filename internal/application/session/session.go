package session

import (
	"bufio"
	"chatY-go/internal/domain/message"
	"chatY-go/internal/domain/user"
	"chatY-go/pkg/logger"
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	CommandNickname = "/nick"
	CommandList     = "/list"
	CommandExit     = "/exit"
	CommandHelp     = "/help"
	CommandWhisper  = "/whisper"
)

type ChatSession struct {
	systemUser user.IUser
	broadcast  chan message.IMessage
	clients    map[net.Conn]user.IUser
	mu         *sync.Mutex
	logger     logger.ILogger
}

func NewChatSession(
	systemUser user.IUser,
	broadcast chan message.IMessage,
	clients map[net.Conn]user.IUser,
	mu *sync.Mutex,
	logger logger.ILogger,
) *ChatSession {
	s := &ChatSession{
		systemUser: systemUser,
		broadcast:  broadcast,
		clients:    clients,
		mu:         mu,
		logger:     logger,
	}
	go s.broadcaster()
	return s
}

func (s *ChatSession) Start(conn net.Conn) {
	defer func() {
		s.logger.Infof("[DISCONNECT] Client %s disconnected", conn.RemoteAddr())
	}()

	s.register(conn)
	s.logger.Infof("[JOIN] User %s joined the chat", s.clients[conn].Nickname())

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
		log.Printf(msg.Print())
		s.logger.Infof("[MESSAGE] %s", msg.Print())

		s.handleMessage(msg, conn)
	}
}

func (s *ChatSession) broadcaster() {
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

func (s *ChatSession) register(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[conn] = user.NewUser(fmt.Sprintf("User_%s", conn.RemoteAddr()))
	s.notify(fmt.Sprintf("User %s joined the chat\n", s.clients[conn].Nickname()))
}

func (s *ChatSession) unregister(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, conn)
	conn.Close()
}
