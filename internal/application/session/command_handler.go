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

func (s *ChatSession) handleListCommand(user user.IUser) {
	users := s.getActiveUsers()
	var b strings.Builder

	b.WriteString("## Active users:\n")
	for _, u := range users {
		b.WriteString(fmt.Sprintf("- %s\n", u.Nickname()))
	}

	s.sendMessageToUser(user, b.String())
}

func (s *ChatSession) handleExitCommand(conn net.Conn) {
	s.notify(fmt.Sprintf("User %s left the chat\n", s.clients[conn].Nickname()))

	s.unregister(conn)
}

func (s *ChatSession) handleHelpCommand(user user.IUser) {
	s.sendMessageToUser(user, HELP_MESSAGE)
}

func (s *ChatSession) handleWhisperCommand(fromUser user.IUser, args []string) {
	if len(args) < 2 {
		s.sendMessageToUser(fromUser, "Usage: /whisper <name> <message>\n")
		return
	}

	toUser := s.getUserByNickname(args[0])
	if toUser == nil {
		s.sendMessageToUser(fromUser, fmt.Sprintf("user %s not found\n", toUser))
		return
	}

	if len(args[1:]) == 0 {
		s.sendMessageToUser(fromUser, "Message are empty, write something\n")
		return
	}

	if fromUser.Nickname() == toUser.Nickname() {
		s.sendMessageToUser(fromUser, "You can't whisper to yourself")
		return
	}

	message := strings.Join(args[1:], " ")
	s.sendMessageToUser(fromUser, fmt.Sprintf("[%s] -> [%s]> %s", fromUser.Nickname(), toUser.Nickname(), message+"\n"))
	s.sendMessageToUser(toUser, fmt.Sprintf("[%s] -> [%s]> %s", fromUser.Nickname(), toUser.Nickname(), message+"\n"))
}
