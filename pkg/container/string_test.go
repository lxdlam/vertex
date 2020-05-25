package container

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/lxdlam/vertex/pkg/protocol"
	"github.com/lxdlam/vertex/pkg/util"
	"github.com/stretchr/testify/assert"
)

const (
	seed                  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	defaultStringLength   = 100
	defaultStringTestCase = 100
)

// Some FNV-1a 64 bits collisions.
// See https://github.com/pstibrany/fnv-1a-64bit-collisions
var collisions = []struct {
	a, b string
}{
	{
		"8yn0iYCKYHlIj4-BwPqk",
		"GReLUrM4wMqfg9yzV3KQ",
	},
	{
		"gMPflVXtwGDXbIhP73TX",
		"LtHf1prlU1bCeYZEdqWf",
	},
	{
		"pFuM83THhM-Qw8FI5FKo",
		".jPx7rOtTDteKAwvfOEo",
	},
	{
		"7mohtcOFVz",
		"c1E51sSEyx",
	},
	{
		"6a5x-VbtXk",
		"f_2k7GG-4v",
	},
}

func TestStringLen(t *testing.T) {
	strs, containers := genRandomCase(defaultListTestCase)

	for idx := 0; idx < defaultStringTestCase; idx++ {
		assert.Equal(t, len(strs[idx]), containers[idx].Len())
	}

	str := string([]byte{97, 98, 99, 100, 0, 0, 101, 102, 103, 49, 50, 51, 52, 53, 54})
	container := NewString(str)

	assert.Equal(t, len(str), container.Len())
}

func TestStringEquals(t *testing.T) {
	// Random
	for idx := 0; idx < defaultListTestCase; idx++ {
		var strA, strB string
		for {
			strA = genRandomString(defaultStringLength)
			strB = genRandomString(defaultStringLength)

			if strA != strB {
				break
			}
		}

		a := NewString(strA)
		b := NewString(strB)

		assert.NotEqual(t, strA, strB)
		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}

	// Hash collisions
	for _, entry := range collisions {
		a := NewString(entry.a)
		b := NewString(entry.b)
		assert.Equal(t, a.Hash(), b.Hash())
		assert.False(t, a.Equals(b))
		assert.False(t, b.Equals(a))
	}

	// Equals same object
	for idx := 0; idx < defaultListTestCase; idx++ {
		obj := NewString(genRandomString(defaultStringLength))

		assert.Equal(t, obj, obj)
		assert.True(t, obj.Equals(obj))
	}

	// Equals different object
	for idx := 0; idx < defaultListTestCase; idx++ {
		str := genRandomString(defaultStringLength)
		a := NewString(str)
		b := NewString(str)

		assert.NotSame(t, a, b)
		assert.True(t, a.Equals(b))
		assert.True(t, b.Equals(a))
	}
}

func TestAppend(t *testing.T) {
	for idx := 0; idx < defaultListTestCase; idx++ {
		var strA, strB string
		for {
			strA = genRandomString(defaultStringLength)
			strB = genRandomString(defaultStringLength)

			if strA != strB {
				break
			}
		}

		a := NewString(strA)
		b := NewString(strB)
		appends := a.Append(b)
		revAppends := b.Append(a)
		aSelfAppends := a.Append(a)
		bSelfAppends := b.Append(b)

		assert.Equal(t, strA, a.String())
		assert.Equal(t, strB, b.String())
		assert.Equal(t, strA+strB, appends.String())
		assert.Equal(t, strB+strA, revAppends.String())
		assert.Equal(t, strA+strA, aSelfAppends.String())
		assert.Equal(t, strB+strB, bSelfAppends.String())
	}
}

func TestStringObject(t *testing.T) {
	// Simple String, only test "OK"
	obj := NewString("OK")
	assert.Equal(t, protocol.NewSimpleRedisString("OK").String(), obj.AsSimpleStringObject().String())

	// Bulk String
	for idx := 0; idx < defaultListTestCase; idx++ {
		str := genRandomString(defaultStringLength)

		assert.Equal(t, protocol.NewBulkRedisString(str).String(), NewString(str).AsBulkStringObject().String())
	}

	// Null Bulk String
	obj = NewString("")
	assert.Equal(t, protocol.NewNullBulkRedisString().String(), obj.AsBulkStringObject().String())
}

func TestStringIntOperations(t *testing.T) {
	for idx := 0; idx < defaultListTestCase; idx++ {
		candidate := util.GetGlobalRandom().Int63()
		if !testIntOperations(t, candidate) {
			t.Fatalf("test string int operations fail. candidate=%d", candidate)
		}
	}

	// Corner cases
	if !testIntOperations(t, 0) {
		t.Fatal("test string int operations fail. candidate=0")
	}
	if !testIntOperations(t, 1) {
		t.Fatal("test string int operations fail. candidate=1")
	}
	if !testIntOperations(t, -1) {
		t.Fatal("test string int operations fail. candidate=-1")
	}
}

func TestRange(t *testing.T) {
	str := "0123456789"
	s := NewString(str)

	for start := 0; start < 10; start++ {
		for end := start; end < 10; end++ {
			actual, err := s.GetRange(start, end)
			assert.Nil(t, err)
			assert.Equal(t, str[start:end+1], actual.String())
		}
	}

	for start := 0; start < 10; start++ {
		for end := start; end < 10; end++ {
			actual, err := s.GetRange(start-10, end)
			assert.Nil(t, err)
			assert.Equal(t, str[start:end+1], actual.String())
		}
	}

	for start := 0; start < 10; start++ {
		for end := start; end < 10; end++ {
			actual, err := s.GetRange(start, end-10)
			assert.Nil(t, err)
			assert.Equal(t, str[start:end+1], actual.String())
		}
	}

	for start := 0; start < 10; start++ {
		for end := start; end < 10; end++ {
			actual, err := s.GetRange(start-10, end-10)
			assert.Nil(t, err)
			assert.Equal(t, str[start:end+1], actual.String())
		}
	}

	_, err := s.GetRange(10, 5)
	assert.EqualError(t, err, "string_container: the given range is invalid")

	_, err = s.GetRange(-1, -2)
	assert.EqualError(t, err, "string_container: the given range is invalid")
}

func TestRangeSpecialCase(t *testing.T) {
	str := NewString("This is a string")

	str1, err := str.GetRange(0, 3)
	assert.Nil(t, err)
	assert.NotNil(t, str1)
	assert.Equal(t, "This", str1.String())

	str2, err := str.GetRange(-3, -1)
	assert.Nil(t, err)
	assert.NotNil(t, str2)
	assert.Equal(t, "ing", str2.String())

	str3, err := str.GetRange(0, -1)
	assert.Nil(t, err)
	assert.NotNil(t, str3)
	assert.Equal(t, "This is a string", str3.String())

	str4, err := str.GetRange(10, 100)
	assert.Nil(t, err)
	assert.NotNil(t, str4)
	assert.Equal(t, "string", str4.String())
}

func testIntOperations(t *testing.T, number int64) bool {
	str := NewString(fmt.Sprintf("%d", number))
	if !testStringEqualInt(t, str, number) {
		t.Errorf("str is not equals to number. str={%+v}, number=%d", str, number)
		return false
	}

	// Increase 1
	ret, err := str.Increase(1)
	if !assert.Nil(t, err) || !assert.Equal(t, number+1, ret) || !testStringEqualInt(t, str, number+1) {
		t.Errorf("str++ failed. str={%+v}, ret=%d, number=%d", str, ret, number)
		return false
	}
	_, _ = str.Decrease(1)

	// Decrease 1
	ret, err = str.Decrease(1)
	if !assert.Nil(t, err) || !assert.Equal(t, number-1, ret) || !testStringEqualInt(t, str, number-1) {
		t.Errorf("str-- failed. str={%+v}, ret=%d, number=%d", str, ret, number)
		return false
	}
	_, _ = str.Increase(1)

	// Increase 0
	ret, err = str.Increase(0)
	if !assert.Nil(t, err) || !assert.Equal(t, number, ret) || !testStringEqualInt(t, str, number) {
		t.Errorf("str += 0 failed. str={%+v}, ret=%d, number=%d", str, ret, number)
		return false
	}

	// Decrease 0
	ret, err = str.Decrease(0)
	if !assert.Nil(t, err) || !assert.Equal(t, number, ret) || !testStringEqualInt(t, str, number) {
		t.Errorf("str -= 0 failed. str={%+v}, ret=%d, number=%d", str, ret, number)
		return false
	}

	// Random increase
	dx := int64(util.GetGlobalRandom().Int31()) // Make sure not exceeded int64
	ret, err = str.Increase(dx)
	if !assert.Nil(t, err) || !assert.Equal(t, number+dx, ret) || !testStringEqualInt(t, str, number+dx) {
		t.Errorf("str += dx failed. str={%+v}, dx=%d, ret=%d, number=%d", str, dx, ret, number)
		return false
	}
	_, _ = str.Decrease(dx)

	// Random decrease
	dx = int64(util.GetGlobalRandom().Int31()) // Make sure not exceeded int64
	ret, err = str.Decrease(dx)
	if !assert.Nil(t, err) || !assert.Equal(t, number-dx, ret) || !testStringEqualInt(t, str, number-dx) {
		t.Errorf("str -= dx failed. str={%+v}, dx=%d, ret=%d, number=%d", str, dx, ret, number)
		return false
	}
	_, _ = str.Increase(dx)

	return true
}

func testStringEqualInt(t *testing.T, str *StringContainer, number int64) bool {
	if !str.IsInt() {
		t.Errorf("str is not an integer. str={%+v}, number=%d", str, number)
		return false
	}

	{
		parsed, err := str.Int()
		if !assert.Equal(t, number, parsed) {
			t.Errorf("the parsed result not equals with number. str={%+v}, number=%d", str, number)
			return false
		}

		if !assert.Nil(t, err) {
			t.Errorf("parse has error. str={%+v}, number=%d", str, number)
			return false
		}
	}

	{
		obj, err := str.AsIntObject()

		if !assert.NotNil(t, obj) || !assert.Nil(t, err) {
			t.Errorf("to redis integer object has error. str={%+v}, number=%d", str, number)
			return false
		}

		if !assert.Equal(t, protocol.NewRedisInteger(number).String(), obj.String()) {
			t.Errorf("to redis integer object not equals the directly built one. str={%+v}, number=%d", str, number)
			return false
		}
	}

	return true
}

func genRandomString(length int) string {
	var buf bytes.Buffer

	for i := 0; i < length; i++ {
		idx := util.GetGlobalRandom().Intn(len(seed))
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
