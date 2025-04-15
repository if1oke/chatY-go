package client

type Model struct {
	Input    string
	Messages []string
	Username string
	Session  IClientSession
}
