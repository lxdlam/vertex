package container

import (
	"errors"
	"hash/fnv"

	"github.com/lxdlam/vertex/pkg/protocol"
)

var (
	dummy = NewString("dummy string")

	// ErrNotAInt will be raised if invoke Int() on a string object that cannot be cast to a string
	ErrNotAInt = errors.New("string_container: cannot cast the value to a string")
)

// StringContainer for containers, which is just a simple type alias
type StringContainer struct {
	data   string
	bytes  []byte
	size   int
	hash   uint64
	intVar intVariant
}

func (s *StringContainer) update(data string) {
	s.data = data
	s.bytes = []byte(data)
	s.size = len(data)

	hash := fnv.New64a()
	_, _ = hash.Write(s.bytes)
	s.hash = hash.Sum64()
}

// NewString will return a pointer to a StringContainer, no copy here
func NewString(s string) *StringContainer {
	str := &StringContainer{}

	str.update(s)
	str.intVar = newIntVariant(s)

	return str
}

// Len will return the size of the string
func (s *StringContainer) Len() int {
	return s.size
}

// Equals will check if the both string is same by their value
func (s *StringContainer) Equals(another *StringContainer) bool {
	if s.hash != another.hash {
		return false
	}

	return s.data == another.data
}

// StringContainer will return the contained string
func (s *StringContainer) String() string {
	return s.data
}

// Byte will return the byte represent of the contained string
func (s *StringContainer) Byte() []byte {
	return s.bytes
}

// Hash will returns the hash value of a string
func (s *StringContainer) Hash() uint64 {
	return s.hash
}

// Append will append another to s and then returns a new string instance
func (s *StringContainer) Append(another *StringContainer) *StringContainer {
	return NewString(s.String() + another.String())
}

// AsSimpleStringObject will return a simple redis string object of the give instance
func (s *StringContainer) AsSimpleStringObject() protocol.RedisString {
	return protocol.NewSimpleRedisString(s.data)
}

// AsBulkStringObject will return a bulk redis string object of the give instance. If the
// string's length is 0, a null bulk redis string will be returned.
func (s *StringContainer) AsBulkStringObject() protocol.RedisString {
	if s.Len() == 0 {
		return protocol.NewNullBulkRedisString()
	}
	return protocol.NewBulkRedisString(s.data)
}

// IsInt reports if the string can be casted into a int
func (s *StringContainer) IsInt() bool {
	return s.intVar != nil
}

// Int will return the go int value of the string
func (s *StringContainer) Int() (int64, error) {
	if s.IsInt() {
		return s.intVar.Get(), nil
	}

	return 0, ErrNotAInt
}

// AsIntObject will return the redis int object of the string
func (s *StringContainer) AsIntObject() (protocol.RedisInteger, error) {
	if s.IsInt() {
		return s.intVar.AsIntObject(), nil
	}

	return nil, ErrNotAInt
}

// Increase will test if the string can be represent as int, then increase the number by the increment.
// The string will be updated too.
func (s *StringContainer) Increase(increment int64) (int64, error) {
	if s.IsInt() {
		s.intVar.Increase(increment)
		s.update(s.intVar.AsString())
		return s.intVar.Get(), nil
	}

	return 0, ErrNotAInt
}

// Decrease will test if the string can be represent as int, then decrease the number by the decrement.
// The string will be updated too.
func (s *StringContainer) Decrease(decrement int64) (int64, error) {
	if s.IsInt() {
		s.intVar.Decrease(decrement)
		s.update(s.intVar.AsString())
		return s.intVar.Get(), nil
	}

	return 0, ErrNotAInt
}

func (s *StringContainer) isContainer() {}

// Key will return the string's value
func (s *StringContainer) Key() string {
	return s.data
}

// Type will return StringType
func (s *StringContainer) Type() ContainerType {
	return StringType
}
