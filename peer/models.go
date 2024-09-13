package peer

import (
	"net"
)

type Message struct {
	Data     []byte
	PeerAddr string
}
type Peer struct {
	Conn    net.Conn
	MsgChan chan Message
	// Memory  memory.Memory
}
