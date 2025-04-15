package client

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Keys events
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			err := m.SendMessage()
			if err != nil {
				m.Messages = append(m.Messages, fmt.Sprintf("SEND ERROR: %s", err))
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			m.Input += msg.String()
		case tea.KeyBackspace:
			m.Input = m.Input[:len(m.Input)-1]
		}

	// Receiver events
	case IncomingMessage:
		m.Messages = append(m.Messages, string(msg))
		return m, nil
	}

	return m, nil
}

func (m *Model) SendMessage() error {
	_, err := m.Session.Writer().WriteString(m.Input + "\n")
	if err != nil {
		return err
	}
	err = m.Session.Writer().Flush()
	if err != nil {
		return err
	}
	m.Input = ""
	return nil
}
