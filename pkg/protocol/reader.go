package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

const (
	delimiter byte = '\n'
)

type respReader struct {
	reader *bufio.Reader
	token  string
}

func (r *respReader) readToken() error {
	var buf bytes.Buffer

	for {
		cur, err := r.reader.ReadString(delimiter)

		if err != nil {
			return fmt.Errorf("read token failed. buf=%s, err={%w}", buf.String(), err)
		}

		buf.WriteString(cur)
		l := len(cur)

		if l <= 1 {
			return fmt.Errorf("read a single '\\n' or empty token. buf=%s", buf.String())
		} else if cur[l-2] == '\r' {
			// We may have abc\ndef\r\n, so we just try to check if we do meet an delimeter
			break
		}
	}

	r.token = buf.String()
	return nil
}

func (r *respReader) readBytes(count int) error {
	var buf bytes.Buffer

	if count < 0 {
		return fmt.Errorf("invalid count. count=%d", count)
	}

	for i := 0; i < count; i++ {
		b, err := r.reader.ReadByte()

		if err != nil {
			return fmt.Errorf("read byte failed. buf=%s, err={%w}", buf.String(), err)
		}

		buf.WriteByte(b)
	}

	r.token = buf.String()
	return nil
}

func (r *respReader) peek() string {
	if len(r.token) <= 0 {
		return ""
	}

	return string(r.token[0])
}

func (r *respReader) readString() (RedisString, error) {
	switch r.peek() {
	case SimpleStringType:
		l := len(r.token)
		return NewSimpleRedisString(r.token[1 : l-2]), nil
	case BulkStringType:
		// length section
		curToken := r.token
		l := len(r.token)
		sLen, err := strconv.ParseInt(curToken[1:l-2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid bulk string length. token=%s, err={%w}", curToken, err)
		}

		// real string section
		if sLen == -1 {
			return NewNullBulkRedisString(), nil
		} else if err := r.readBytes(int(sLen) + 2); err == nil {
			return NewBulkRedisString(r.token[:sLen]), nil
		} else {
			return nil, fmt.Errorf("read real string failed. length=%d, token=%s, err={%w}", sLen, r.token, err)
		}
	}

	return nil, fmt.Errorf("unknown string type. token=%s", r.token)
}

func (r *respReader) readError() (RedisError, error) {
	l := len(r.token)
	return NewRedisError(r.token[1 : l-2]), nil
}

func (r *respReader) readInteger() (RedisInteger, error) {
	l := len(r.token)
	num, err := strconv.ParseInt(r.token[1:l-2], 10, 64)

	if err != nil {
		return nil, err
	}
	return NewRedisInteger(num), nil
}

func (r *respReader) readArray() (RedisArray, error) {
	// length section
	curToken := r.token
	l := len(r.token)
	aLen, err := strconv.ParseInt(curToken[1:l-2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid array length. token=%s, err={%w}", curToken, err)
	}

	if aLen == -1 {
		// null array
		return NewNullRedisArray(), nil
	}

	var objects []RedisObject
	for i := 0; i < int(aLen); i++ {
		obj, err := r.readObject()
		if err != nil {
			return nil, fmt.Errorf("read array object failed. index=%d, err={%w}", i, err)
		}
		objects = append(objects, obj)
	}

	return NewRedisArray(objects), nil
}

func (r *respReader) readObject() (RedisObject, error) {
	if err := r.readToken(); err == nil {
		switch r.peek() {
		case SimpleStringType, BulkStringType:
			return r.readString()
		case ErrorType:
			return r.readError()
		case IntegerType:
			return r.readInteger()
		case ArrayType:
			return r.readArray()
		default:
			return nil, fmt.Errorf("invalid token. token=%s", r.token)
		}
	} else {
		return nil, fmt.Errorf("peek front meet error. token=%s, err={%w}", r.token, err)
	}
}

// Parse takes an io.Reader and try to parse a RedisObject from it
//
// If any error raised, a wrapped ErrInvalidRESP error will be returned
func Parse(reader io.Reader) (RedisObject, error) {
	r := &respReader{
		reader: bufio.NewReader(reader),
	}

	return r.readObject()
}
