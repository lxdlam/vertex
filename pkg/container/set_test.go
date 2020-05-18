package container

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

const (
	defaultSetTestCase = 100
)

func TestSetBasicOperation(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	s := NewSetContainer("test")

	_, items := genRandomCase(defaultSetTestCase)

	s.Add(items)

	rand.Shuffle(defaultSetTestCase, func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	// IsMember
	for _, item := range items {
		assert.True(t, s.IsMember(item))
	}

	// Delete
	for _, item := range items {
		s.Delete([]*StringContainer{item})
		assert.False(t, s.IsMember(item))
	}

	rand.Shuffle(defaultSetTestCase, func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	for _, item := range items {
		assert.False(t, s.IsMember(item))
	}
}

func TestDuplicateItems(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	s := NewSetContainer("test")
	var entries []*StringContainer

	for idx := 0; idx < defaultSetTestCase; idx++ {
		entries = append(entries, NewString(fmt.Sprintf("%d", rand.Int63()%5)))
	}

	s.Add(entries)

	assert.Equal(t, 5, s.Len())

	for item := 0; item < 5; item++ {
		assert.True(t, s.IsMember(NewString(fmt.Sprintf("%d", item))))
	}
}

func TestSetCollisionEntries(t *testing.T) {
	for _, entry := range collisions {
		s := NewSetContainer("test")
		entries := []*StringContainer{NewString(entry.a), NewString(entry.b)}

		s.Add(entries)

		assert.Equal(t, 2, s.Len())

		for _, entry := range entries {
			assert.True(t, s.IsMember(entry))
		}

		members := s.Members()

		assert.ElementsMatch(t, entries, members)
	}
}

func TestSetAccess(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	s := NewSetContainer("test")

	_, items := genRandomCase(defaultSetTestCase)

	s.Add(items)

	rand.Shuffle(defaultSetTestCase, func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	members := s.Members()
	assert.ElementsMatch(t, items, members)
}

func TestSetPop(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	s := NewSetContainer("test")

	_, items := genRandomCase(defaultSetTestCase)

	s.Add(items)

	rand.Shuffle(defaultSetTestCase, func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	randomMembers := s.RandomMember(defaultSetTestCase + 20)
	assert.Equal(t, defaultSetTestCase, s.Len())
	assert.ElementsMatch(t, items, randomMembers)

	// Negative number is acceptable
	randomMembers = s.RandomMember(-defaultSetTestCase)
	assert.Equal(t, defaultSetTestCase, s.Len())
	assert.ElementsMatch(t, items, randomMembers)

	popedMembers := s.Pop(-defaultSetTestCase - 20)
	assert.Equal(t, 0, s.Len())
	assert.ElementsMatch(t, items, popedMembers)

	assert.Nil(t, s.RandomMember(100))
	assert.Nil(t, s.Pop(100))
}

// TODO: We may need more test case?
func TestSetDiff(t *testing.T) {
	// Redis case
	a, b, c := genRedisTestCase()
	expected := []*StringContainer{NewString("b"), NewString("d")}
	actual := a.Diff([]SetContainer{b, c}).Members()
	assert.ElementsMatch(t, expected, actual)
}

// TODO: We may need more test case?
func TestSetIntersect(t *testing.T) {
	// Redis case
	a, b, c := genRedisTestCase()
	expected := []*StringContainer{NewString("c")}
	actual := a.Intersect([]SetContainer{b, c}).Members()
	assert.ElementsMatch(t, expected, actual)
}

// TODO: We may need more test case?
func TestSetUnion(t *testing.T) {
	a, b, c := genRedisTestCase()
	expected := []*StringContainer{NewString("a"), NewString("b"), NewString("c"), NewString("d"), NewString("e")}
	actual := a.Union([]SetContainer{b, c}).Members()
	assert.ElementsMatch(t, expected, actual)
}

func genRedisTestCase() (SetContainer, SetContainer, SetContainer) {
	a := NewSetContainer("a")
	a.Add([]*StringContainer{
		NewString("a"),
		NewString("b"),
		NewString("c"),
		NewString("d"),
	})

	b := NewSetContainer("b")
	b.Add([]*StringContainer{
		NewString("c"),
	})

	c := NewSetContainer("c")
	c.Add([]*StringContainer{
		NewString("a"),
		NewString("c"),
		NewString("e"),
	})

	return a, b, c
}
