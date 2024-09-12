package server

import (
	"log"
	"net"

	"github.com/muhammadolammi/tinydb/aof"
	"github.com/muhammadolammi/tinydb/memory"
)

type Server struct {
	ListerAddr string
	ln         net.Listener
	MsgChan    chan []byte
	AOFMsgChan chan []byte

	Conn   net.Conn
	Memory *memory.Memory
	Aof    *aof.AOF

	// peers       map[net.Addr]*peer.Peer
	// addPeerChan chan *peer.Peer
	quitChan chan struct{}
}

func NewServer(listerAddr string) *Server {
	aof, err := aof.NewAOF("data/database.aof")
	if err != nil {
		log.Panicf("error creating aof file . error: %s", err)
	}
	return &Server{
		ListerAddr: listerAddr,
		// peers:       make(map[net.Addr]*peer.Peer),
		// addPeerChan: make(chan *peer.Peer),
		quitChan:   make(chan struct{}),
		MsgChan:    make(chan []byte),
		AOFMsgChan: make(chan []byte),
		Memory:     memory.NewMemory(),
		Aof:        aof,
	}
}
