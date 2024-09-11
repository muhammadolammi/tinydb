package server

import (
	"bytes"
	"io"
	"log/slog"
	"net"
)

func NewServer(listerAddr string) *Server {
	return &Server{
		ListerAddr: listerAddr,
		// peers:       make(map[net.Addr]*peer.Peer),
		// addPeerChan: make(chan *peer.Peer),
		quitChan: make(chan struct{}),
		MsgChan:  make(chan []byte),
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

		case rawMsg := <-s.MsgChan:
			s.handleRawMessage(bytes.NewReader(rawMsg))
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
		s.Conn = conn
		go s.readLoop()
	}
}

func (s *Server) readLoop() {
	buff := make([]byte, 1024)
	for {
		n, err := s.Conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			slog.Info("read loop error.", "error", err)

		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buff)
		// fmt.Println(msgBuf)
		s.MsgChan <- msgBuf
	}
}

// func (s *Server) handleConn() {
// 	peer := peer.NewPeer(conn, s.MsgChan)
// 	s.addPeerChan <- peer
// 	slog.Info("new peer added", "remoteAddr", conn.RemoteAddr())

// 	peer.ReadLoop()
// 	if err != nil {
// 		slog.Info("read loop error.", "error", err, "remoteAddr", s.Conn.RemoteAddr())
// 	}
// }

func (s *Server) handleRawMessage(r io.Reader) {
	s.Conn.Write([]byte("+OK\r\n"))

}
