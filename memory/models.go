package memory

import (
	"sync"

	"github.com/muhammadolammi/tinydb/protocol"
)

type Memory struct {
	mu       sync.RWMutex
	Hash     map[string]string
	Handlers map[string]func(args []protocol.Value) *protocol.Value
	hsmu     sync.RWMutex
	HHash    map[string]map[string]string
}

func NewMemory() *Memory {
	m := &Memory{
		Hash:  make(map[string]string),
		HHash: make(map[string]map[string]string),
	}

	m.Handlers = map[string]func(args []protocol.Value) *protocol.Value{
		"PING":    m.ping,
		"GET":     m.get,
		"SET":     m.set,
		"HSET":    m.hset,
		"HGET":    m.hget,
		"HGETALL": m.hgetall,
	}

	return m
}
