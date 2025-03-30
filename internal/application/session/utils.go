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
		s.logger.Errorf("[ERROR] Write to client %s failed: %v", conn.RemoteAddr(), err)
	}
}

func (s *ChatSession) notify(text string) {
	s.doBroadcast(message.NewMessage(s.systemUser, text))
}

func (s *ChatSession) doBroadcast(message message.IMessage) {
	s.broadcast <- message
}
