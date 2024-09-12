package protocol

import (
	"bufio"
	"io"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

type RESP struct {
	reader *bufio.Reader
}

type RESPWriter struct {
	writer io.Writer
}

func NewResp(rd io.Reader) *RESP {
	return &RESP{reader: bufio.NewReader(rd)}
}

func NewRespWriter(w io.Writer) *RESPWriter {
	return &RESPWriter{writer: w}
}
