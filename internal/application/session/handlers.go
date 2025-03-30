package session

import (
	"chatY-go/internal/domain/message"
	"chatY-go/internal/domain/user"
	"fmt"
	"net"
	"strings"
)

const (
	HELP_MESSAGE = "Available commands:\n/nick <name> - change nickname\n/list - show active users\n/exit - disconnect\n/help - this help message\n"
)

func (s *ChatSession) handleMessage(message message.IMessage, conn net.Conn) {
	cmd, arg := parseCommand(message.Text())

	switch cmd {
	case CommandNickname:
		s.handleNicknameCommand(message.User(), arg[0])
	case CommandList:
		s.handleListCommand(message.User())
	case CommandExit:
		s.handleExitCommand(conn)
	case CommandHelp:
		s.handleHelpCommand(message.User())
	case CommandWhisper:
		s.handleWhisperCommand(message.User(), arg)
	default:
		s.doBroadcast(message)
	}
}

func (s *ChatSession) handleNicknameCommand(usr user.IUser, arg string) {
	if arg == "" {
		s.notify("Usage: /nick <new_nickname>\n")
		return
	}

	old := usr.Nickname()

	s.mu.Lock()
	usr.SetNickname(arg)
	s.mu.Unlock()

	s.logger.Infof("[NICK] %s changed nickname to %s", old, arg)
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
	s.logger.Infof("[LEAVE] User %s left the chat", s.clients[conn].Nickname())

	s.unregister(conn)
}

func (s *ChatSession) handleHelpCommand(user user.IUser) {
	s.sendMessageToUser(user, HELP_MESSAGE)
}

func (s *ChatSession) handleWhisperCommand(fromUser user.IUser, args []string) {
	if len(args) < 2 {
		s.sendMessageToUser(fromUser, "Usage: /whisper <name> <message>\n")
		s.logger.Warnf("[WHISPER FAIL] wrong command format %s", args)
		return
	}

	toUser := s.getUserByNickname(args[0])
	if toUser == nil {
		s.sendMessageToUser(fromUser, fmt.Sprintf("user %s not found\n", toUser))
		s.logger.Warnf("[WHISPER FAIL] user %s not found", args[0])
		return
	}

	if fromUser.Nickname() == toUser.Nickname() {
		s.sendMessageToUser(fromUser, "You can't whisper to yourself")
		return
	}

	message := strings.Join(args[1:], " ")

	if strings.TrimSpace(message) == "" {
		s.sendMessageToUser(fromUser, "Message are empty, write something\n")
		s.logger.Warn("[WHISPER FAIL] message are empty")
		return
	}

	s.logger.Infof("[WHISPER] %s -> %s: %s", fromUser.Nickname(), toUser.Nickname(), message)

	s.sendMessageToUser(fromUser, fmt.Sprintf("[%s] -> [%s]> %s", fromUser.Nickname(), toUser.Nickname(), message+"\n"))
	s.sendMessageToUser(toUser, fmt.Sprintf("[%s] -> [%s]> %s", fromUser.Nickname(), toUser.Nickname(), message+"\n"))
}
