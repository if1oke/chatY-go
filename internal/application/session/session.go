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

		s.handleCommands(msg)

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

func (s *ChatSession) handleCommands(message message.IMessage) {
	msgArr := strings.Split(strings.Replace(message.Text(), "\n", "", 1), " ")
	switch msgArr[0] {
	case "/nick":
		oldNickname := message.User().Nickname()
		s.setNickname(message.User(), msgArr[1])
		message.SetText(fmt.Sprintf("%s nickname changed to %s\n", oldNickname, msgArr[1]))
	}
}

func (s *ChatSession) doBroadcast(message message.IMessage) {
	s.broadcast <- message
}

func (s *ChatSession) register(conn net.Conn) {
	s.mu.Lock()
	s.clients[conn] = user.NewUser(fmt.Sprintf("User_%s", conn.RemoteAddr()))
	s.mu.Unlock()
}

func (s *ChatSession) unregister(conn net.Conn) {
	s.mu.Lock()
	delete(s.clients, conn)
	s.mu.Unlock()
}

func (s *ChatSession) setNickname(user user.IUser, nickname string) {
	s.mu.Lock()
	user.SetNickname(nickname)
	s.mu.Unlock()
}
