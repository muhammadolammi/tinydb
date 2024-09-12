package protocol

import (
	"strconv"
)

func (r *RESP) Read() (*Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return &Value{}, err
	}
	switch _type {
	case BULK:
		return r.readBulk()
	case ARRAY:
		return r.readArray()
	default:
		return &Value{}, UnkownCommandError

	}

	// return &Value{}, nil
}
func (r *RESP) readBulk() (*Value, error) {
	value := Value{}
	value.Typ = "bulk"
	// this is just a bulk, so we just read the line and save it to value.

	// lets get the lenght of the bulk to know how long we read
	bulkLen, _, err := r.readInteger()
	if err != nil {
		return &value, err
	}
	bulk := make([]byte, bulkLen)
	r.reader.Read(bulk)

	value.Bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return &value, nil
}
func (r *RESP) readArray() (*Value, error) {
	value := Value{}
	value.Typ = "array"

	// lets get the lenght of the array to know how long we read
	arrayLen, _, err := r.readInteger()
	if err != nil {
		return &value, err
	}

	// v.array = make([]Value, 0)
	// for i := 0; i < len; i++ {
	// 	val, err := r.Read()
	// 	if err != nil {
	// 		return v, err
	// 	}

	// 	// append parsed value to array
	// 	v.array = append(v.array, val)
	// }
	// in this array implementation, we need to make some decision
	// if an error happens underway in any index should we return the captured value
	// or just return nil
	// i went for the fact that we must get all value or nil
	// for now
	array := make([]Value, arrayLen)
	for i := 0; i < arrayLen; i++ {

		v, err := r.Read()
		if err != nil {
			return &value, err
		}
		array[i] = *v
	}
	value.Array = array

	return &value, nil
}

func (r *RESP) readLine() ([]byte, int, error) {
	line := []byte{}
	n := 0
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return line, n, err
		}
		line = append(line, b)
		n++
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}

	}

	return line[:len(line)-2], n, nil

}

func (r *RESP) readInteger() (int, int, error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(i64), n, err
}
