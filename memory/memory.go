package memory

import (
	"github.com/muhammadolammi/tinydb/protocol"
)

func (m *Memory) set(args []protocol.Value) *protocol.Value {
	if len(args) != 2 {
		return &protocol.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}

	}
	m.mu.Lock()
	m.Hash[args[0].Bulk] = args[1].Bulk
	m.mu.Unlock()

	return &protocol.Value{Typ: "string", Str: "OK"}
}

func (m *Memory) get(args []protocol.Value) *protocol.Value {
	if len(args) != 1 {
		return &protocol.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}

	}

	key := args[0].Bulk

	m.mu.RLock()
	value, ok := m.Hash[key]
	m.mu.RUnlock()

	if !ok {
		return &protocol.Value{Typ: "null"}
	}

	return &protocol.Value{Typ: "bulk", Bulk: value}
}

func (m *Memory) hset(args []protocol.Value) *protocol.Value {
	if len(args) != 3 {
		return &protocol.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	m.hsmu.Lock()
	if _, ok := m.HHash[hash]; !ok {
		m.HHash[hash] = map[string]string{}
	}
	m.HHash[hash][key] = value
	m.hsmu.Unlock()

	return &protocol.Value{Typ: "string", Str: "OK"}
}

func (m *Memory) hget(args []protocol.Value) *protocol.Value {
	if len(args) != 2 {
		return &protocol.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	m.hsmu.RLock()
	value, ok := m.HHash[hash][key]
	m.hsmu.RUnlock()

	if !ok {
		return &protocol.Value{Typ: "null"}
	}

	return &protocol.Value{Typ: "bulk", Bulk: value}
}

func (m *Memory) hgetall(args []protocol.Value) *protocol.Value {
	if len(args) != 1 {
		return &protocol.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	m.hsmu.RLock()
	values, ok := m.HHash[hash]
	m.hsmu.RUnlock()

	if !ok {
		return &protocol.Value{Typ: "null"}
	}
	arr := make([]protocol.Value, 0)
	for key, val := range values {
		arr = append(arr, protocol.Value{Typ: "string", Str: key})
		arr = append(arr, protocol.Value{Typ: "string", Str: val})

	}

	return &protocol.Value{Typ: "array", Array: arr}
}
func (m *Memory) ping(args []protocol.Value) *protocol.Value {
	// if args is empty return with pong
	// else return with the first string of the args
	if len(args) == 0 {
		return &protocol.Value{Typ: "string", Str: "PONG"}
	}
	return &protocol.Value{Typ: "string", Str: args[0].Bulk}
}
