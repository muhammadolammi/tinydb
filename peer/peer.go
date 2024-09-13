package peer

import (
	"net"
)

func NewPeer(conn net.Conn, msgChan chan Message) *Peer {
	return &Peer{
		Conn:    conn,
		MsgChan: msgChan,
	}
}

func (p *Peer) ReadLoop() error {
	buff := make([]byte, 1024)
	for {
		n, err := p.Conn.Read(buff)
		if err != nil {
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buff)
		// fmt.Println(msgBuf)
		message := Message{
			Data:     msgBuf,
			PeerAddr: p.Conn.RemoteAddr().String(),
		}
		p.MsgChan <- message
	}
}
