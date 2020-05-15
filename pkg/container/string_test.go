package container_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/lxdlam/vertex/pkg/container"
)

func testStringEqual(t *testing.T, lhs, rhs *StringContainer) bool {
	if !assert.Equal(t, lhs.String(), rhs.String()) {
		return false
	}

	if !assert.ElementsMatch(t, lhs.Byte(), rhs.Byte()) {
		return false
	}

	if !assert.Equal(t, lhs.Len(), rhs.Len()) {
		return false
	}

	if !assert.Equal(t, lhs.Hash(), rhs.Hash()) {
		return false
	}

	if !assert.Equal(t, lhs.IsInt(), rhs.IsInt()) {
		return false
	}

	lval, lerr := lhs.Int()
	rval, rerr := rhs.Int()

	if !assert.Equal(t, lval, rval) || !assert.Equal(t, lerr, rerr) {
		return false
	}

	return true
}
