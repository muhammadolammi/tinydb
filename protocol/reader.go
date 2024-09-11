package protocol

import (
	"strconv"
)

func (r *RESP) Read() error {
	b, _ := r.reader.ReadByte()
	if b != BULK {
		return WrongCommandType

	}
	// read next  byte for the size of the array
	size, _ := r.reader.ReadByte()

	strSize, _ := strconv.ParseInt(string(size), 10, 64)

	// consume /r/n
	r.reader.ReadByte()
	r.reader.ReadByte()

	return nil
}
