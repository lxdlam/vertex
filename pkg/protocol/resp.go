package protocol

import (
	"bytes"
	"errors"
	"fmt"
)

// We use the first byte to indicate the object type
const (
	SimpleStringType = "+"
	ErrorType        = "-"
	IntegerType      = ":"
	BulkStringType   = "$"
	ArrayType        = "*"

	Delimiter = "\r\n"

	NullBulkStringLiteral = "$-1\r\n"
	NullArrayLiteral      = "*-1\r\n"
)

// RedisObject is the RESP object interface.
type RedisObject interface {
	// Byte will return a byte slice which can be directly write to the io.Writer
	Byte() []byte

	// RedisString will return the string representation of the inner byte buffer.
	String() string

	// Type will return the RedisObject's type, which can help you to check before you cast.
	Type() string
}

// RedisString is both the RESP simple string and bulk string interface.
type RedisString interface {
	RedisObject

	// Data returns the inner data.
	Data() string
}

// RedisError is the RESP error interface.
type RedisError interface {
	RedisObject

	// Message returns a string, indicates the error message.
	Message() string

	// Error returns a go error, which message is set to the real message.
	Error() error
}

// RedisInteger is the RESP integer interface.
type RedisInteger interface {
	RedisObject

	// Data returns an integer. It's set to int64, cast to the suitable representation.
	Data() int64
}

// RedisArray is the RESP RedisArray interface.
type RedisArray interface {
	RedisObject

	// Data returns an RedisObject slice include the inner value in the array
	Data() []RedisObject
}

type redisString struct {
	data      string
	stringRep string
	byteRep   []byte
	subType   string
}

// NewSimpleRedisString takes an string, return a new RedisString instance which Type()==SimpleStringType
func NewSimpleRedisString(data string) RedisString {
	s := fmt.Sprintf("%s%s%s", SimpleStringType, data, Delimiter)
	return &redisString{
		data:      data,
		stringRep: s,
		byteRep:   []byte(s),
		subType:   SimpleStringType,
	}
}

// NewBulkRedisString takes an string, return a new RedisString instance which Type()==BulkStringType
func NewBulkRedisString(data string) RedisString {
	s := fmt.Sprintf("%s%d%s%s%s", BulkStringType, len(data), Delimiter, data, Delimiter)
	return &redisString{
		data:      data,
		stringRep: s,
		byteRep:   []byte(s),
		subType:   BulkStringType,
	}
}

// NewNullBulkRedisString will return a NullBulkString, i.e., "$-1\r\n"
func NewNullBulkRedisString() RedisString {
	return &redisString{
		data:      "",
		stringRep: NullBulkStringLiteral,
		byteRep:   []byte(NullBulkStringLiteral),
		subType:   BulkStringType,
	}
}

func (rs *redisString) Byte() []byte {
	return rs.byteRep
}

func (rs *redisString) String() string {
	return rs.stringRep
}

func (rs *redisString) Type() string {
	return rs.subType
}

func (rs *redisString) Data() string {
	return rs.data
}

type redisError struct {
	data      string
	stringRep string
	byteRep   []byte
}

// NewRedisError takes an string, return a new RedisError instance
func NewRedisError(data string) RedisError {
	s := fmt.Sprintf("%s%s%s", ErrorType, data, Delimiter)

	return &redisError{
		data:      data,
		stringRep: s,
		byteRep:   []byte(s),
	}
}

// NewRedisErrorFromGoError takes an go error, return a new RedisError instance
func NewRedisErrorFromGoError(data error) RedisError {
	s := fmt.Sprintf("%s%s%s", ErrorType, data.Error(), Delimiter)

	return &redisError{
		data:      data.Error(),
		stringRep: s,
		byteRep:   []byte(s),
	}
}

func (re *redisError) Byte() []byte {
	return re.byteRep
}

func (re *redisError) String() string {
	return re.stringRep
}

func (re *redisError) Type() string {
	return ErrorType
}

func (re *redisError) Message() string {
	return re.data
}

func (re *redisError) Error() error {
	return errors.New(re.data)
}

type redisInteger struct {
	data      int64
	stringRep string
	byteRep   []byte
}

// NewRedisInteger takes an int64, return a new RedisInteger instance
func NewRedisInteger(data int64) RedisInteger {
	s := fmt.Sprintf("%s%d%s", IntegerType, data, Delimiter)

	return &redisInteger{
		data:      data,
		stringRep: s,
		byteRep:   []byte(s),
	}
}

func (ri *redisInteger) Byte() []byte {
	return ri.byteRep
}

func (ri *redisInteger) String() string {
	return ri.stringRep
}

func (ri *redisInteger) Type() string {
	return IntegerType
}

func (ri *redisInteger) Data() int64 {
	return ri.data
}

type redisArray struct {
	data      []RedisObject
	stringRep string
	byteRep   []byte
}

// NewRedisArray takes a RedisObject slice, return a new RedisArray instance
func NewRedisArray(data []RedisObject) RedisArray {
	var buf bytes.Buffer

	ra := &redisArray{}

	buf.WriteString(fmt.Sprintf("%s%d%s", ArrayType, len(data), Delimiter))

	for _, obj := range data {
		ra.data = append(ra.data, obj)
		buf.Write(obj.Byte())
	}

	ra.byteRep = buf.Bytes()
	ra.stringRep = buf.String()

	return ra
}

// NewNullRedisArray will return a NullRedisArray, i.e., "*-1\r\n"
func NewNullRedisArray() RedisArray {
	return &redisArray{
		data:      nil,
		stringRep: NullArrayLiteral,
		byteRep:   []byte(NullArrayLiteral),
	}
}

func (ra *redisArray) Byte() []byte {
	return ra.byteRep
}

func (ra *redisArray) String() string {
	return ra.stringRep
}

func (ra *redisArray) Type() string {
	return ArrayType
}

func (ra *redisArray) Data() []RedisObject {
	return ra.data
}
