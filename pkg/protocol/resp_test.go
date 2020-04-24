package protocol_test

import (
	"errors"
	"testing"

	. "github.com/lxdlam/vertex/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestSimpleString(t *testing.T) {
	obj := NewSimpleRedisString("OK")

	assert.Equal(t, "+OK\r\n", obj.String())
}

func TestError(t *testing.T) {
	obj := NewRedisError("Error Message")

	assert.Equal(t, "-Error Message\r\n", obj.String())
}

func TestErrorFromGo(t *testing.T) {
	obj := NewRedisErrorFromGoError(errors.New("Go Error"))

	assert.Equal(t, "-Go Error\r\n", obj.String())
}

func TestInteger(t *testing.T) {
	obj := NewRedisInteger(1000)

	assert.Equal(t, ":1000\r\n", obj.String())
}

func TestZero(t *testing.T) {
	obj := NewRedisInteger(0)

	assert.Equal(t, ":0\r\n", obj.String())
}

func TestNegative(t *testing.T) {
	obj := NewRedisInteger(-10000)

	assert.Equal(t, ":-10000\r\n", obj.String())
}

func TestEmptyBulkString(t *testing.T) {
	obj := NewBulkRedisString("")

	assert.Equal(t, "$0\r\n\r\n", obj.String())
}

func TestBulkString(t *testing.T) {
	// A string mixed with \r and \n
	obj := NewBulkRedisString("Hello \rWorld!\n")

	assert.Equal(t, "$14\r\nHello \rWorld!\n\r\n", obj.String())
}

func TestComplexBulkString(t *testing.T) {
	obj := NewBulkRedisString("$198\r\nMcpJxQSaoHhgLbAsdPTRjGFCZqrwlEivnIuKkfOVYzWDNXmyUetB\r\n\r\n\r\x0c \n\t\x0b\r\n\r\n^|,!%#`:_<*-=?;$[\\&.>\"(~+]'/}{)@\r\n\r\nETbgCoRcXOWezIHKxmQqVvhADyGZJplufNFjakdYPrwMLsBitUSn\r\n\r\n_)~&!\\.[<=*]^>+;$%(/'@,}\"{:-|#`?\r\n\r\n\r\n \t\r\n")

	assert.Equal(t, "$206\r\n$198\r\nMcpJxQSaoHhgLbAsdPTRjGFCZqrwlEivnIuKkfOVYzWDNXmyUetB\r\n\r\n\r\x0c \n\t\x0b\r\n\r\n^|,!%#`:_<*-=?;$[\\&.>\"(~+]'/}{)@\r\n\r\nETbgCoRcXOWezIHKxmQqVvhADyGZJplufNFjakdYPrwMLsBitUSn\r\n\r\n_)~&!\\.[<=*]^>+;$%(/'@,}\"{:-|#`?\r\n\r\n\r\n \t\r\n\r\n", obj.String())
}

func TestEmptyArray(t *testing.T) {
	obj := NewRedisArray(nil)

	assert.Equal(t, "*0\r\n", obj.String())
}

func TestArray(t *testing.T) {
	obj := NewRedisArray([]RedisObject{
		NewSimpleRedisString("OK"),
		NewRedisError("Error Message"),
		NewRedisInteger(1000),
		NewBulkRedisString("Hello \rWorld!\n"),
	})

	assert.Equal(t, "*4\r\n+OK\r\n-Error Message\r\n:1000\r\n$14\r\nHello \rWorld!\n\r\n", obj.String())
}

func TestNullInArray(t *testing.T) {
	obj := NewRedisArray([]RedisObject{
		NewBulkRedisString("foo"),
		NewNullBulkRedisString(),
		NewBulkRedisString("bar"),
	})

	assert.Equal(t, "*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n", obj.String())
}

func TestNestedArray(t *testing.T) {
	obj := NewRedisArray([]RedisObject{
		NewRedisArray([]RedisObject{
			NewRedisArray([]RedisObject{
				NewRedisArray([]RedisObject{
					NewRedisInteger(-1),
					NewSimpleRedisString("OK"),
				}),
			}),
		}),
	})

	assert.Equal(t, "*1\r\n*1\r\n*1\r\n*2\r\n:-1\r\n+OK\r\n", obj.String())
}

func TestComplexArray(t *testing.T) {
	// Mix above all
	obj := NewRedisArray([]RedisObject{
		NewSimpleRedisString("OK"),
		NewRedisError("Error Message"),
		NewRedisErrorFromGoError(errors.New("Go Error")),
		NewRedisInteger(1000),
		NewRedisInteger(0),
		NewRedisInteger(-10000),
		NewBulkRedisString(""),
		NewBulkRedisString("Hello \rWorld!\n"),
		NewBulkRedisString("$198\r\nMcpJxQSaoHhgLbAsdPTRjGFCZqrwlEivnIuKkfOVYzWDNXmyUetB\r\n\r\n\r\x0c \n\t\x0b\r\n\r\n^|,!%#`:_<*-=?;$[\\&.>\"(~+]'/}{)@\r\n\r\nETbgCoRcXOWezIHKxmQqVvhADyGZJplufNFjakdYPrwMLsBitUSn\r\n\r\n_)~&!\\.[<=*]^>+;$%(/'@,}\"{:-|#`?\r\n\r\n\r\n \t\r\n"),
		NewRedisArray(nil),
		NewRedisArray([]RedisObject{
			NewSimpleRedisString("OK"),
			NewRedisError("Error Message"),
			NewRedisInteger(1000),
			NewBulkRedisString("Hello \rWorld!\n"),
		}),
		NewRedisArray([]RedisObject{
			NewBulkRedisString("foo"),
			NewNullBulkRedisString(),
			NewBulkRedisString("bar"),
		}),
		NewRedisArray([]RedisObject{
			NewRedisArray([]RedisObject{
				NewRedisArray([]RedisObject{
					NewRedisArray([]RedisObject{
						NewRedisInteger(-1),
						NewSimpleRedisString("OK"),
					}),
				}),
			}),
		}),
		// Why not have a null array? ;)
		NewNullRedisArray(),
	})

	assert.Equal(t, "*14\r\n+OK\r\n-Error Message\r\n-Go Error\r\n:1000\r\n:0\r\n:-10000\r\n$0\r\n\r\n$14\r\nHello \rWorld!\n\r\n$206\r\n$198\r\nMcpJxQSaoHhgLbAsdPTRjGFCZqrwlEivnIuKkfOVYzWDNXmyUetB\r\n\r\n\r\x0c \n\t\x0b\r\n\r\n^|,!%#`:_<*-=?;$[\\&.>\"(~+]'/}{)@\r\n\r\nETbgCoRcXOWezIHKxmQqVvhADyGZJplufNFjakdYPrwMLsBitUSn\r\n\r\n_)~&!\\.[<=*]^>+;$%(/'@,}\"{:-|#`?\r\n\r\n\r\n \t\r\n\r\n*0\r\n*4\r\n+OK\r\n-Error Message\r\n:1000\r\n$14\r\nHello \rWorld!\n\r\n*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n*1\r\n*1\r\n*1\r\n*2\r\n:-1\r\n+OK\r\n*-1\r\n", obj.String())
}
