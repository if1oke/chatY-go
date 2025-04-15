package session

import (
	"bufio"
	"net"
)

type ClientSession struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	username string
	password string
}

func (s *ClientSession) Reader() *bufio.Reader {
	return s.reader
}

func (s *ClientSession) Writer() *bufio.Writer {
	return s.writer
}

func NewSession(addr string, username, password string) (*ClientSession, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &ClientSession{
		conn:     conn,
		reader:   bufio.NewReader(conn),
		writer:   bufio.NewWriter(conn),
		username: username,
		password: password,
	}, nil
}
