package session

import (
	"bufio"
	"chatY-go/internal/domain/message"
	"chatY-go/internal/domain/user"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	COMMAND_NICKNAME = "/nick"
	COMMAND_LIST     = "/list"
	COMMAND_EXIT     = "/exit"
)

type ChatSession struct {
	broadcast chan message.IMessage
	clients   map[net.Conn]user.IUser
	mu        *sync.Mutex
}

func NewChatSession(
	broadcast chan message.IMessage,
	clients map[net.Conn]user.IUser,
	mu *sync.Mutex,
) *ChatSession {
	s := &ChatSession{
		broadcast: broadcast,
		clients:   clients,
		mu:        mu,
	}
	go s.broadcaster()
	return s
}

func (s *ChatSession) Start(conn net.Conn) {
	s.register(conn)

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
			log.Printf("Client disconnected: %v", err)
			s.unregister(conn)
			return
		}

		msg := message.NewMessage(s.clients[conn], rawMessage)
		log.Printf(msg.Print())

		s.handleCommands(msg, conn)

		s.doBroadcast(msg)
	}
}

func (s *ChatSession) broadcaster() {
	for {
		msg := <-s.broadcast
		for client := range s.clients {
			_, err := fmt.Fprintf(client, msg.Print())
			if err != nil {
				log.Printf("write to client error: %s", err.Error())
			}
		}
	}
}

func (s *ChatSession) handleCommands(message message.IMessage, conn net.Conn) {
	cmd, arg := parseCommand(message.Text())

	switch cmd {
	case COMMAND_NICKNAME:
		s.handleNicknameCommand(message, arg)
	case COMMAND_LIST:
		s.handleListCommand(message)
	case COMMAND_EXIT:
		s.unregister(conn)
	}
}

func (s *ChatSession) doBroadcast(message message.IMessage) {
	s.broadcast <- message
}

func (s *ChatSession) register(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[conn] = user.NewUser(fmt.Sprintf("User_%s", conn.RemoteAddr()))
}

func (s *ChatSession) unregister(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, conn)
	conn.Close()
}

func (s *ChatSession) getActiveUsers() []user.IUser {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []user.IUser

	for _, v := range s.clients {
		users = append(users, v)
	}

	return users
}

func parseCommand(text string) (string, string) {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\r", ""))
	parts := strings.SplitN(text, " ", 2)

	cmd := parts[0]
	arg := ""
	if len(parts) > 1 {
		arg = strings.TrimSpace(parts[1])
	}

	return cmd, arg
}
