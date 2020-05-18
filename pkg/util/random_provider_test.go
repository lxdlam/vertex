package util_test

import (
	. "github.com/lxdlam/vertex/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNewRandom(t *testing.T) {
	rand1 := GetNewRandom()
	rand2 := GetNewRandom()

	assert.NotNil(t, rand1)
	assert.NotNil(t, rand2)
	assert.NotSame(t, rand1, rand2)

	for i := 0; i < 100; i++ {
		assert.NotEqual(t, rand1.Int63(), rand2.Int63())
	}
}

func TestGetGlobalRandom(t *testing.T) {
	rand1 := GetGlobalRandom()
	rand2 := GetGlobalRandom()

	assert.NotNil(t, rand1)
	assert.NotNil(t, rand2)
	assert.Same(t, rand1, rand2)
}
