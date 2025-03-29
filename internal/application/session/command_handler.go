package session

import (
	"chatY-go/internal/domain/user"
	"fmt"
	"net"
	"strings"
)

const (
	HELP_MESSAGE = "Available commands:\n/nick <name> - change nickname\n/list - show active users\n/exit - disconnect\n/help - this help message\n"
)

func (s *ChatSession) handleNicknameCommand(usr user.IUser, arg string) {
	if arg == "" {
		s.notify("Usage: /nick <new_nickname>\n")
		return
	}

	old := usr.Nickname()

	s.mu.Lock()
	usr.SetNickname(arg)
	s.mu.Unlock()

	s.notify(fmt.Sprintf("%s nickname changed to %s\n", old, arg))
}

func (s *ChatSession) handleListCommand() {
	users := s.getActiveUsers()
	var b strings.Builder

	b.WriteString("## Active users:\n")
	for _, u := range users {
		b.WriteString(fmt.Sprintf("- %s\n", u.Nickname()))
	}

	s.notify(b.String())
}

func (s *ChatSession) handleExitCommand(conn net.Conn) {
	s.notify(fmt.Sprintf("User %s left the chat\n", s.clients[conn].Nickname()))

	s.unregister(conn)
}

func (s *ChatSession) handleHelpCommand() {
	s.notify(HELP_MESSAGE)
}
