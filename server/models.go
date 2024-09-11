package server

import (
	"net"
)

type Server struct {
	ListerAddr string
	ln         net.Listener
	MsgChan    chan []byte
	Conn       net.Conn

	// peers       map[net.Addr]*peer.Peer
	// addPeerChan chan *peer.Peer
	quitChan chan struct{}
}
