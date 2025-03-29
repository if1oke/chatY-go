package session

import (
	"chatY-go/internal/domain/message"
	"fmt"
	"strings"
)

func (s *ChatSession) handleNicknameCommand(message message.IMessage, arg string) {
	if arg == "" {
		message.SetSystemText("Usage: /nick <new_nickname>\n")
		return
	}

	old := message.User().Nickname()

	s.mu.Lock()
	message.User().SetNickname(arg)
	s.mu.Unlock()

	message.SetSystemText(fmt.Sprintf("%s nickname changed to %s\n", old, arg))
}

func (s *ChatSession) handleListCommand(message message.IMessage) {
	users := s.getActiveUsers()
	var b strings.Builder

	b.WriteString("## Active users:\n")
	for _, u := range users {
		b.WriteString(fmt.Sprintf("- %s\n", u.Nickname()))
	}

	message.SetSystemText(b.String())
}
