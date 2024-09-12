package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"strings"

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

		case rawMsg := <-s.MsgChan:
			s.handleRawMessage(bytes.NewReader(rawMsg))
		case aofrawMsg := <-s.AOFMsgChan:
			s.handleAOFRawMessage(bytes.NewReader(aofrawMsg))
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

func (s *Server) handleConn(conn net.Conn) {
	// peer := peer.NewPeer(conn, s.MsgChan)
	// s.addPeerChan <- peer
	// slog.Info("new peer added", "remoteAddr", conn.RemoteAddr())

	// peer.ReadLoop()
	// if err != nil {
	// 	slog.Info("read loop error.", "error", err, "remoteAddr", s.Conn.RemoteAddr())
	// }

	// lets read from the loop
	s.Conn = conn
	// go s.readAOFAndSendToChannel()
	s.readLoop()

	defer s.Conn.Close()
}

func (s *Server) handleRawMessage(r io.Reader) {
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

	writer := protocol.NewRespWriter(s.Conn)

	handler, ok := s.Memory.Handlers[command]
	if !ok {
		errString := fmt.Sprintf("Invalid command: %s", command)
		writer.Write(&protocol.Value{Typ: "error", Str: errString})
		return

	}
	// lets save commands to memory
	if command == "SET" || command == "HSET" {
		s.Aof.Write(*value)
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

	writer := protocol.NewRespWriter(s.Conn)

	handler, ok := s.Memory.Handlers[command]
	if !ok {
		errString := fmt.Sprintf("Invalid command: %s", command)
		writer.Write(&protocol.Value{Typ: "error", Str: errString})
		return

	}
	// // lets save commands to memory
	// if command == "SET" || command == "HSET" {
	// 	s.Aof.Write(*value)
	// }
	_ = handler(args)
	// writer.Write(result)

}
