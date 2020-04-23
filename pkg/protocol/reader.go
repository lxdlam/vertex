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

func (r *respReader) nextToken() error {
	var buf bytes.Buffer

	for {
		cur, err := r.reader.ReadString('\n')

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

func (r *respReader) peekFront() (string, error) {
	if err := r.nextToken(); err != nil {
		return "", err
	}

	if len(r.token) <= 0 {
		return "", fmt.Errorf("empty token")
	}

	return string(r.token[0]), nil
}

func (r *respReader) readString() (RedisString, error) {
	switch string(r.token[0]) {
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

		// null bulk string
		if sLen == -1 {
			return NewNullBulkRedisString(), nil
		}

		// real string section
		if err := r.nextToken(); err != nil {
			return nil, fmt.Errorf("invalid bulk string second token. token=%s, err={%w}", curToken, err)
		}

		if sLen == 0 {
			return NewBulkRedisString(""), nil
		} else if int(sLen) == len(r.token)-2 {
			return NewBulkRedisString(r.token[1 : l-2]), nil
		}

		return nil, fmt.Errorf("bulk string length mismatch. length=%d, token=%s", sLen, r.token)
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
	if front, err := r.peekFront(); err == nil {
		switch front {
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
