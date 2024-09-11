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
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type RESP struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *RESP {
	return &RESP{reader: bufio.NewReader(rd)}
}
