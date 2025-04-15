package client

import (
	tea "github.com/charmbracelet/bubbletea"
)

type IApplication interface {
	Config() IClientConfig
	Session() IClientSession
}

type IncomingMessage string

func (m *Model) Init() tea.Cmd {
	return nil
}

func Run(app IApplication) error {
	m := Model{
		Username: app.Config().Username(),
		Input:    "",
		Messages: []string{},
		Session:  app.Session(),
	}

	err := m.Session.Authenticate()
	if err != nil {
		return err
	}

	p := tea.NewProgram(&m)

	go func() {
		for {
			msg, err := app.Session().Reader().ReadString('\n')
			if err != nil {
				break
			}
			p.Send(IncomingMessage(msg))
		}
	}()

	_, err = p.Run()

	return err
}
