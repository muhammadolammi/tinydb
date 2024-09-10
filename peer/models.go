package peer

import "net"

type Peer struct {
	Conn    net.Conn
	MsgChan chan []byte
}
