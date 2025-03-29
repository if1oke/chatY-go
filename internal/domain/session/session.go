package session

import "net"

type ISession interface {
	Start(conn net.Conn)
}
