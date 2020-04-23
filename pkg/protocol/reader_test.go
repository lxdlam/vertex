package protocol_test

import (
	"strings"
	"testing"

	. "github.com/lxdlam/vertex/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestParseEmptyBulkString(t *testing.T) {
	raw := "$0\r\n\r\n"

	obj, err := Parse(strings.NewReader(raw))

	assert.Nil(t, err)
	assert.Equal(t, "$0\r\n\r\n", obj.String())
}

func TestParseEmptyArray(t *testing.T) {
	raw := "*0\r\n"

	obj, err := Parse(strings.NewReader(raw))

	assert.Nil(t, err)
	assert.Equal(t, "*0\r\n", obj.String())
}

func TestParseNullBulkString(t *testing.T) {
	raw := "$-1\r\n"

	obj, err := Parse(strings.NewReader(raw))

	assert.Nil(t, err)
	assert.Equal(t, "$-1\r\n", obj.String())
}

func TestParseNullArray(t *testing.T) {
	raw := "*-1\r\n"

	obj, err := Parse(strings.NewReader(raw))

	assert.Nil(t, err)
	assert.Equal(t, "*-1\r\n", obj.String())
}
