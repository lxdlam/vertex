package protocol_test

import (
	"testing"

	. "github.com/lxdlam/vertex/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestEmptyArray(t *testing.T) {
	obj := NewRedisArray(nil)

	assert.Equal(t, "*0\r\n", obj.String())
}

func TestEmptyBulkString(t *testing.T) {
	obj := NewBulkRedisString("")

	assert.Equal(t, "$0\r\n\r\n", obj.String())
}
