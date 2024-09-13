package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"strings"

	"github.com/muhammadolammi/tinydb/peer"
	"github.com/muhammadolammi/tinydb/protocol"
)

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListerAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.readAOFAndSendToChannel()
	go s.loop()
	slog.Info("server running", "listenAddr", s.ListerAddr)
	defer s.Aof.Close()
	return s.acceptLoop()
}

func (s *Server) loop() {
	for {
		select {

		case msg := <-s.MsgChan:
			p := s.peers[msg.PeerAddr]
			s.handleRawMessage(bytes.NewReader(msg.Data), *p)
		case aofrawMsg := <-s.AOFMsgChan:
			s.handleAOFRawMessage(bytes.NewReader(aofrawMsg))
		case peer := <-s.addPeerChan:
			s.peers[peer.Conn.RemoteAddr().String()] = peer
			slog.Info("new peer added", "remoteAddr", peer.Conn.RemoteAddr())
		case peer := <-s.removePeerChan:
			delete(s.peers, peer.Conn.RemoteAddr().String())
			slog.Info("peer removed", "remoteAddr", peer.Conn.RemoteAddr())

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

// func (s *Server) readLoop() {
// 	buff := make([]byte, 1024)
// 	for {
// 		n, err := s.Conn.Read(buff)
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			slog.Info("read loop error.", "error", err)

// 		}
// 		msgBuf := make([]byte, n)
// 		copy(msgBuf, buff)
// 		// fmt.Println(msgBuf)
// 		s.MsgChan <- msgBuf
// 	}
// }

func (s *Server) handleConn(conn net.Conn) {
	peer := peer.NewPeer(conn, s.MsgChan)
	s.addPeerChan <- peer

	err := peer.ReadLoop()
	if err != nil {
		if err == io.EOF {
			s.removePeerChan <- peer
			return
		}

		slog.Info("read loop error.", "error", err, "remoteAddr", peer.Conn.RemoteAddr())
	}

	// lets read from the loop
	defer peer.Conn.Close()
}

func (s *Server) handleRawMessage(r io.Reader, peer peer.Peer) {
	resp := protocol.NewResp(r)

	// check if the reader is an array or not
	value, err := resp.Read()

	if err != nil {
		log.Println("raw message error", "err", err)
		s.quitChan <- struct{}{}
	}
	if value.Typ != "array" {
		log.Println("Invalid request, expected array")

	}

	if len(value.Array) == 0 {
		log.Println("Invalid request, expected array length > 0")
	}
	command := strings.ToUpper(value.Array[0].Bulk)
	args := value.Array[1:]
	// fmt.Println(args)

	writer := protocol.NewRespWriter(peer.Conn)

	handler, ok := s.Memory.Handlers[command]
	if !ok {
		errString := fmt.Sprintf("Invalid command: %s", command)
		writer.Write(&protocol.Value{Typ: "error", Str: errString})
		return

	}
	// lets save commands to memory
	if command == "SET" || command == "HSET" {
		err = s.Aof.Write(*value)
		if err != nil {
			log.Println("error writing to aof", "error", err)
		}
	}
	result := handler(args)
	writer.Write(result)

}

func (s *Server) readAOFAndSendToChannel() {
	for {
		resp := protocol.NewResp(s.Aof.Reader)
		v, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				// EOF reached, exit the loop

				return
			}
			// Log the error and continue reading
			log.Printf("Error reading AOF file: %v", err)
			continue
		}
		b := v.Marshal()

		s.AOFMsgChan <- b
	}
}

func (s *Server) handleAOFRawMessage(r io.Reader) {
	resp := protocol.NewResp(r)

	// check if the reader is an array or not
	value, err := resp.Read()

	if err != nil {
		log.Println("raw message error", "err", err)
		return
	}
	if value.Typ != "array" {
		log.Println("Invalid request, expected array")

		return

	}

	if len(value.Array) == 0 {
		log.Println("Invalid request, expected array length > 0")
		return
	}
	command := strings.ToUpper(value.Array[0].Bulk)
	args := value.Array[1:]
	// fmt.Println(args)

	handler, ok := s.Memory.Handlers[command]
	if !ok {
		errString := fmt.Sprintf("Invalid command: %s", command)
		// writer.Write(&protocol.Value{Typ: "error", Str: errString})
		log.Println(errString)
		return

	}

	_ = handler(args)
}
