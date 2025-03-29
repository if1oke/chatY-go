package session

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type ChatSession struct {
	broadcast chan string
	clients   map[net.Conn]bool
	mu        *sync.Mutex
}

func NewChatSession(
	broadcast chan string,
	clients map[net.Conn]bool,
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
	s.Register(conn)

	defer func() {
		s.Unregister(conn)
		err := conn.Close()
		if err != nil {
			return
		}
	}()

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client disconnected: %v", err)
			s.Unregister(conn)
			return
		}
		log.Printf("Received: %s", msg)

		s.Broadcast(msg)
	}
}

func (s *ChatSession) broadcaster() {
	for {
		msg := <-s.broadcast
		for client := range s.clients {
			_, err := fmt.Fprintf(client, msg)
			if err != nil {
				log.Printf("write to client error: %s", err.Error())
			}
		}
	}
}

func (s *ChatSession) Broadcast(message string) {
	s.broadcast <- message
}

func (s *ChatSession) Register(conn net.Conn) {
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()
}

func (s *ChatSession) Unregister(conn net.Conn) {
	s.mu.Lock()
	delete(s.clients, conn)
	s.mu.Unlock()
}
