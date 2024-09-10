package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/muhammadolammi/tinydb/peer"
)

func NewServer(listerAddr string) *Server {
	return &Server{
		ListerAddr:  listerAddr,
		peers:       make(map[net.Addr]*peer.Peer),
		addPeerChan: make(chan *peer.Peer),
		quitChan:    make(chan struct{}),
		MsgChan:     make(chan []byte),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListerAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()
	slog.Info("server running", "listenAddr", s.ListerAddr)
	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerChan:
			s.peers[s.ln.Addr()] = peer
		case rawMsg := <-s.MsgChan:
			fmt.Println(rawMsg)
		case <-s.quitChan:
			return

			// default:
			// 	time.Sleep(10 * time.Second)
			// 	log.Println("default peer looping")

		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Info("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := peer.NewPeer(conn, s.MsgChan)
	s.addPeerChan <- peer
	slog.Info("new peer added", "remoteAddr", conn.RemoteAddr())

	err := peer.ReadLoop()
	if err != nil {
		slog.Info("peer read error.", "error", err, "remoteAddr", conn.RemoteAddr())
	}
}
