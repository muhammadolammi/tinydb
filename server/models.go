package server

import (
	"net"

	"github.com/muhammadolammi/tinydb/peer"
)

type Server struct {
	ListerAddr string
	ln         net.Listener
	MsgChan    chan []byte

	peers       map[net.Addr]*peer.Peer
	addPeerChan chan *peer.Peer
	quitChan    chan struct{}
}
