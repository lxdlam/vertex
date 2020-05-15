package container

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	seed                = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultStringLength = 100
	defaultTestCase     = 100
)

func TestPushAndPopHead(t *testing.T) {
	l := NewLinkedListContainer("test")

	expected, testCase := genRandomCase(defaultTestCase)
	expected = reverse(expected)

	size, err := l.PushHead(testCase)

	assert.Nil(t, err)
	assert.Equal(t, len(expected), size)
	assert.ElementsMatch(t, expected, l.(*linkedList).debugExtract())

	var actual []string

	for idx := len(expected) - 1; idx >= 0; idx-- {
		ret, err := l.PopHead()
		assert.Nil(t, err)
		assert.Equal(t, idx, l.Len())
		actual = append(actual, ret.data)
	}

	assert.Equal(t, 0, l.Len())
	assert.ElementsMatch(t, expected, actual)
}

func TestPushAndPopTail(t *testing.T) {
	l := NewLinkedListContainer("test")

	expected, testCase := genRandomCase(defaultTestCase)

	size, err := l.PushTail(testCase)

	assert.Nil(t, err)
	assert.Equal(t, len(expected), size)
	assert.ElementsMatch(t, expected, l.(*linkedList).debugExtract())

	var actual []string

	for idx := len(expected) - 1; idx >= 0; idx-- {
		ret, err := l.PopTail()
		assert.Nil(t, err)
		assert.Equal(t, idx, l.Len())
		actual = append(actual, ret.data)
	}

	assert.Equal(t, 0, l.Len())
	assert.ElementsMatch(t, expected, actual)
}

func genRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	var buf bytes.Buffer

	for i := 0; i < length; i++ {
		idx := rand.Intn(len(seed))
		buf.WriteByte(seed[idx])
	}

	return buf.String()
}

func genRandomCase(size int) ([]string, []*StringContainer) {
	var expected []string
	var testCase []*StringContainer

	for i := 0; i < size; i++ {
		str := genRandomString(defaultStringLength)
		expected = append(expected, str)
		testCase = append(testCase, NewString(str))
	}

	return expected, testCase
}

func reverse(slice []string) []string {
	ret := slice[:]

	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}

	return ret
}
