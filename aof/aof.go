package aof

import (
	"github.com/muhammadolammi/tinydb/protocol"
)

func (aof *AOF) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *AOF) Write(value protocol.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	// fmt.Println(string(value.Marshal()))

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}
