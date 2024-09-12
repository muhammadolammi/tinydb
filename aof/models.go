package aof

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type AOF struct {
	file   *os.File
	Reader *bufio.Reader
	mu     sync.Mutex
}

func NewAOF(filePath string) (*AOF, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return &AOF{}, err
	}
	aof := &AOF{
		file:   file,
		Reader: bufio.NewReader(file),
	}
	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}
