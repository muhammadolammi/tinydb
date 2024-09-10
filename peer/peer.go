package peer

import (
	"net"
)

func NewPeer(conn net.Conn, msgChan chan []byte) *Peer {
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
		p.MsgChan <- msgBuf
	}
}
