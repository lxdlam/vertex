package container

import (
	"fmt"
	"testing"

	"github.com/lxdlam/vertex/pkg/util"

	"github.com/stretchr/testify/assert"
)

const (
	defaultHashTestCase = 100
)

func TestHashBasicOperation(t *testing.T) {
	h := NewHashContainer("test")

	vis := make(map[string]bool)

	_, keys := genRandomCase(defaultHashTestCase)
	_, values := genRandomCase(defaultHashTestCase)

	var randomKeys []*StringContainer

	for _, item := range keys {
		vis[item.String()] = true
	}

	for idx := 0; idx < defaultHashTestCase; idx++ {
		for {
			str := genRandomString(defaultStringLength)
			if _, ok := vis[str]; !ok {
				randomKeys = append(randomKeys, NewString(str))
				break
			}
		}
	}

	err := h.Set(keys, values)
	assert.Nil(t, err)

	// Get
	getValues := h.Get(keys)
	assert.Equal(t, len(keys), len(getValues))
	for idx := 0; idx < defaultHashTestCase; idx++ {
		// In the same order of key
		assert.NotNil(t, getValues[idx])
		assert.Same(t, values[idx], getValues[idx])
	}

	getRandomValues := h.Get(randomKeys)
	assert.Equal(t, len(randomKeys), len(getRandomValues))
	for idx := 0; idx < defaultHashTestCase; idx++ {
		assert.Nil(t, getRandomValues[idx])
	}

	// Exists
	for idx := 0; idx < defaultHashTestCase; idx++ {
		assert.True(t, h.Exists(keys[idx]))
		assert.False(t, h.Exists(randomKeys[idx]))
	}

	// Del
	for idx := 0; idx < defaultHashTestCase; idx++ {
		h.Del(keys[idx])
		ret := h.Get([]*StringContainer{keys[idx]})
		assert.Equal(t, 1, len(ret))
		assert.Nil(t, ret[0])
	}

	getDeletedValues := h.Get(keys)
	assert.Equal(t, len(keys), len(getDeletedValues))
	for idx := 0; idx < defaultHashTestCase; idx++ {
		assert.Nil(t, getDeletedValues[idx])
	}

	for idx := 0; idx < defaultHashTestCase; idx++ {
		h.Del(randomKeys[idx])
		ret := h.Get([]*StringContainer{randomKeys[idx]})
		assert.Equal(t, 1, len(ret))
		assert.Nil(t, ret[0])
	}

	getRandomValues = h.Get(randomKeys)
	assert.Equal(t, len(randomKeys), len(getRandomValues))
	for idx := 0; idx < defaultHashTestCase; idx++ {
		assert.Nil(t, getRandomValues[idx])
	}

	// Set Errors
	err = h.Set(keys, values[1:])
	assert.Equal(t, ErrHashLengthNotMatch, err)
}

func TestHashCollisionKeys(t *testing.T) {
	values := []*StringContainer{NewString("a"), NewString("b")}

	for _, entry := range collisions {
		h := NewHashContainer("test")
		keys := []*StringContainer{NewString(entry.a), NewString(entry.b)}

		err := h.Set(keys, values)
		assert.Nil(t, err)

		assert.Equal(t, 2, h.Len())
		getValues := h.Get(keys)
		assert.ElementsMatch(t, values, getValues)
	}
}

func TestHashIncrease(t *testing.T) {
	h := NewHashContainer("test")
	keys := []*StringContainer{NewString("string"), NewString("int")}
	values := []*StringContainer{NewString("abc"), NewString("1")}

	err := h.Set(keys, values)
	assert.Nil(t, err)

	// Increase int
	dx := util.GetGlobalRandom().Int63()
	ret, err := h.Increase(keys[1], dx)

	getValues := h.Get([]*StringContainer{keys[1]})
	assert.Equal(t, 1, len(getValues))
	assert.Same(t, getValues[0], values[1])

	assert.Nil(t, err)
	assert.Equal(t, 1+dx, ret)
	assert.True(t, testStringEqualInt(t, values[1], 1+dx))

	// Increase a string
	ret, err = h.Increase(keys[0], dx)
	assert.Equal(t, int64(0), ret)
	assert.EqualError(t, err, "hash_container: increase with error. err={string_container: cannot cast the value to a string}")

	// Increase a non-exists key
	ret, err = h.Increase(NewString("foo"), dx)
	assert.Equal(t, int64(0), ret)
	assert.Equal(t, ErrKeyNotExist, err)
}

func TestHashExtract(t *testing.T) {
	h := NewHashContainer("test")

	_, keys := genRandomCase(defaultHashTestCase)
	_, values := genRandomCase(defaultHashTestCase)

	err := h.Set(keys, values)
	assert.Nil(t, err)

	extractKeys := h.Keys()
	assert.ElementsMatch(t, keys, extractKeys)

	extractValues := h.Values()
	assert.ElementsMatch(t, values, extractValues)

	// Iterate all entries to check if both are match
	extractKeys, extractValues = h.Entries()
	for idx := 0; idx < defaultHashTestCase; idx++ {
		for j := 0; j < defaultHashTestCase; j++ {
			if extractKeys[idx] == keys[j] {
				assert.Same(t, values[j], extractValues[idx])
			}
		}
	}
}

func TestHashKeyLen(t *testing.T) {
	h := NewHashContainer("test")

	var keys, values []*StringContainer

	for idx := 0; idx < defaultHashTestCase; idx++ {
		keys = append(keys, NewString(fmt.Sprintf("%d", idx+1)))
		values = append(values, NewString(genRandomString(idx+1)))
	}

	err := h.Set(keys, values)
	assert.Nil(t, err)

	util.GetGlobalRandom().Shuffle(defaultHashTestCase, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for _, key := range keys {
		expected, _ := key.Int()
		actual, _ := h.KeyLen(key)

		assert.Equal(t, int(expected), actual)
	}

	// Non-exist keys
	length, err := h.KeyLen(NewString("abc"))
	assert.Equal(t, 0, length)
	assert.Equal(t, ErrKeyNotExist, err)
}
