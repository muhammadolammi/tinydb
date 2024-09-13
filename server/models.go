package server

import (
	"log"
	"net"

	"github.com/muhammadolammi/tinydb/aof"
	"github.com/muhammadolammi/tinydb/memory"
	"github.com/muhammadolammi/tinydb/peer"
)

type Server struct {
	ListerAddr string
	ln         net.Listener
	MsgChan    chan peer.Message
	AOFMsgChan chan []byte

	Memory *memory.Memory
	Aof    *aof.AOF

	peers          map[string]*peer.Peer
	addPeerChan    chan *peer.Peer
	removePeerChan chan *peer.Peer

	quitChan chan struct{}
}

func NewServer(listerAddr string) *Server {
	aof, err := aof.NewAOF("data/database.aof")
	if err != nil {
		log.Panicf("error creating aof file . error: %s", err)
	}
	return &Server{
		ListerAddr:     listerAddr,
		peers:          make(map[string]*peer.Peer),
		addPeerChan:    make(chan *peer.Peer),
		removePeerChan: make(chan *peer.Peer),

		quitChan: make(chan struct{}),

		MsgChan:    make(chan peer.Message),
		AOFMsgChan: make(chan []byte),
		Memory:     memory.NewMemory(),
		Aof:        aof,
	}
}
