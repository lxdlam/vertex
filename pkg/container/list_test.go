package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	defaultListTestCase = 100
)

func TestListPushAndPopHead(t *testing.T) {
	l := NewLinkedListContainer("test")

	expected, testCase := genRandomCase(defaultListTestCase)
	expected = reverse(expected)

	size, err := l.PushHead(testCase)

	assert.Nil(t, err)
	assert.Equal(t, len(expected), size)
	assert.ElementsMatch(t, expected, extractRange(l, 0, -1))

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

func TestListPushAndPopTail(t *testing.T) {
	l := NewLinkedListContainer("test")

	expected, testCase := genRandomCase(defaultListTestCase)

	size, err := l.PushTail(testCase)

	assert.Nil(t, err)
	assert.Equal(t, len(expected), size)
	assert.ElementsMatch(t, expected, extractRange(l, 0, -1))

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

func TestListInsert(t *testing.T) {
	l := NewLinkedListContainer("test")

	expected, testCase := genRandomCase(defaultListTestCase)

	// Ensure in the same order
	size, err := l.PushTail(testCase)

	assert.Nil(t, err)
	assert.Equal(t, len(expected), size)
	assert.ElementsMatch(t, expected, extractRange(l, 0, -1))

}

func TestListSet(t *testing.T) {

}

func TestListRemove(t *testing.T) {

}

func TestListTrim(t *testing.T) {

}

func TestListIndex(t *testing.T) {

}

func TestListRange(t *testing.T) {
	l := NewLinkedListContainer("test")

	expected, testCase := genRandomCase(defaultListTestCase)

	// Ensure in the same order
	size, err := l.PushTail(testCase)

	assert.Nil(t, err)
	assert.Equal(t, len(expected), size)
	assert.ElementsMatch(t, expected, extractRange(l, 0, -1))


}

func extractRange(l ListContainer, left, right int) []string {
	var ret []string
	list, _ := l.Range(left, right)
	for _, item := range list {
		ret = append(ret, item.data)
	}

	return ret
}
