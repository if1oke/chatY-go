package session

import (
	"chatY-go/internal/domain/message"
	"chatY-go/internal/domain/user"
	"fmt"
	"net"
	"strings"
)

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

func (s *ChatServer) getActiveUsers() []user.IUser {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users []user.IUser

	for _, v := range s.clients {
		users = append(users, v)
	}

	return users
}

func (s *ChatServer) getUserByNickname(nickname string) user.IUser {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.clients {
		if v.Nickname() == nickname {
			return v
		}
	}
	return nil
}

func (s *ChatServer) getConnByNickname(nickname string) net.Conn {
	s.mu.Lock()
	defer s.mu.Unlock()

	for c, v := range s.clients {
		if v.Nickname() == nickname {
			return c
		}
	}
	return nil
}

func (s *ChatServer) sendMessageToUser(user user.IUser, text string) {
	conn := s.getConnByNickname(user.Nickname())
	_, err := fmt.Fprintf(conn, text)
	if err != nil {
		s.logger.Errorf("[ERROR] Write to client_old %s failed: %v", conn.RemoteAddr(), err)
	}
}

func (s *ChatServer) sendMessageToConn(conn net.Conn, text string) {
	_, err := fmt.Fprintf(conn, text)
	if err != nil {
		s.logger.Errorf("[ERROR] Write to client_old %s failed: %v", conn.RemoteAddr(), err)
	}
}

func (s *ChatServer) notify(text string) {
	s.doBroadcast(message.NewMessage(s.systemUser, text))
}

func (s *ChatServer) doBroadcast(message message.IMessage) {
	s.broadcast <- message
}
