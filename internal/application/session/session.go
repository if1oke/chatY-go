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
	COMMAND_HELP     = "/help"
	COMMAND_WHISPER  = "/whisper"
)

type ChatSession struct {
	systemUser user.IUser
	broadcast  chan message.IMessage
	clients    map[net.Conn]user.IUser
	mu         *sync.Mutex
}

func NewChatSession(
	systemUser user.IUser,
	broadcast chan message.IMessage,
	clients map[net.Conn]user.IUser,
	mu *sync.Mutex,
) *ChatSession {
	s := &ChatSession{
		systemUser: systemUser,
		broadcast:  broadcast,
		clients:    clients,
		mu:         mu,
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

		s.handleMessage(msg, conn)
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

func (s *ChatSession) handleMessage(message message.IMessage, conn net.Conn) {
	cmd, arg := parseCommand(message.Text())

	switch cmd {
	case COMMAND_NICKNAME:
		s.handleNicknameCommand(message.User(), arg[0])
	case COMMAND_LIST:
		s.handleListCommand(message.User())
	case COMMAND_EXIT:
		s.handleExitCommand(conn)
	case COMMAND_HELP:
		s.handleHelpCommand(message.User())
	case COMMAND_WHISPER:
		s.handleWhisperCommand(message.User(), arg)
	default:
		s.doBroadcast(message)
	}
}

func (s *ChatSession) notify(text string) {
	s.doBroadcast(message.NewMessage(s.systemUser, text))
}

func (s *ChatSession) doBroadcast(message message.IMessage) {
	s.broadcast <- message
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

func (s *ChatSession) getActiveUsers() []user.IUser {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []user.IUser

	for _, v := range s.clients {
		users = append(users, v)
	}

	return users
}

func (s *ChatSession) getUserByNickname(nickname string) user.IUser {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.clients {
		if v.Nickname() == nickname {
			return v
		}
	}
	return nil
}

func (s *ChatSession) getConnByNickname(nickname string) net.Conn {
	s.mu.Lock()
	defer s.mu.Unlock()

	for c, v := range s.clients {
		if v.Nickname() == nickname {
			return c
		}
	}
	return nil
}

func (s *ChatSession) sendMessageToUser(user user.IUser, text string) {
	conn := s.getConnByNickname(user.Nickname())
	_, err := fmt.Fprintf(conn, text)
	if err != nil {
		log.Printf("write to client error: %s", err.Error())
	}
}

func parseCommand(text string) (string, []string) {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\r", ""))
	parts := strings.Split(text, " ")

	cmd := parts[0]
	args := parts[1:]
	if len(parts) > 1 {
		for i, v := range args {
			args[i] = strings.TrimSpace(v)
		}
	}

	return cmd, args
}
