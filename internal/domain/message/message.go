package message

import (
	"chatY-go/internal/domain/user"
	"fmt"
)

type IMessage interface {
	User() user.IUser
	Text() string
	SetText(text string)
	SystemText() string
	SetSystemText(text string)
	Print() string
}

type Message struct {
	user       user.IUser
	text       string
	systemText string
}

func (m *Message) SystemText() string {
	return m.systemText
}

func (m *Message) SetSystemText(systemText string) {
	m.systemText = systemText
}

func NewMessage(user user.IUser, text string) *Message {
	return &Message{user: user, text: text}
}

func (m *Message) User() user.IUser {
	return m.user
}

func (m *Message) Text() string {
	return m.text
}

func (m *Message) SetText(text string) {
	m.text = text
}

func (m *Message) Print() string {
	message := fmt.Sprintf("[%s]> %s", m.user.Nickname(), m.text)
	if m.SystemText() != "" {
		message += fmt.Sprintf("\n%s", m.SystemText())
	}
	return message
}
