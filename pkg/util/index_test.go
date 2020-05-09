package util_test

import (
	"testing"

	. "github.com/lxdlam/vertex/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestIndexResolveRaw(t *testing.T) {
	testCases := []struct {
		index    int
		size     int
		expected int
	}{
		{
			0,
			10,
			0,
		},
		{
			1,
			10,
			1,
		},
		{
			10,
			10,
			10,
		},
		{
			20,
			10,
			20,
		},
		{
			-1,
			10,
			9,
		},
		{
			-10,
			10,
			0,
		},
		{
			-11,
			10,
			-1,
		},
		{
			0,
			0,
			0,
		},
		{
			1,
			0,
			1,
		},
		{
			-1,
			0,
			-1,
		},
	}

	for idx, testCase := range testCases {
		actual := NewIndex(testCase.index).ResolveRaw(testCase.size)
		if !assert.Equal(t, testCase.expected, actual) {
			t.Fatalf("ResolveRaw result differ. idx=%d, expected=%d, actual=%d", idx, testCase.expected, actual)
		}
	}
}

func TestIndexResolve(t *testing.T) {
	testCases := []struct {
		index    int
		size     int
		expected int
	}{
		{
			0,
			10,
			0,
		},
		{
			1,
			10,
			1,
		},
		{
			10,
			10,
			-1,
		},
		{
			20,
			10,
			-1,
		},
		{
			-1,
			10,
			9,
		},
		{
			-10,
			10,
			0,
		},
		{
			-11,
			10,
			-1,
		},
		{
			0,
			0,
			-1,
		},
		{
			1,
			0,
			-1,
		},
		{
			-1,
			0,
			-1,
		},
	}

	for idx, testCase := range testCases {
		actual := NewIndex(testCase.index).Resolve(testCase.size)
		if !assert.Equal(t, testCase.expected, actual) {
			t.Fatalf("ResolveRaw result differ. idx=%d, expected=%d, actual=%d", idx, testCase.expected, actual)
		}
	}
}
