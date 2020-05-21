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

// Assume a single point can resolved correctly, we just need to check if both can resolved correctly
// and if the segment itself is invalid can be detected correctly.
func TestSliceResolveResolve(t *testing.T) {
	testCases := []struct {
		indexLeft     int
		indexRight    int
		size          int
		expectedLeft  int
		expectedRight int
	}{
		// In the range
		{
			3,
			9,
			10,
			3,
			9,
		},
		{
			3,
			-4,
			10,
			3,
			6,
		},
		{
			-7,
			9,
			10,
			3,
			9,
		},
		{
			0,
			-1,
			10,
			0,
			9,
		},
		{
			-6,
			-4,
			10,
			4,
			6,
		},
		// Single Point
		{
			7,
			7,
			10,
			7,
			7,
		},
		// Out of range
		{
			3,
			10,
			10,
			-1,
			-1,
		},
		{
			20,
			30,
			10,
			-1,
			-1,
		},
		{
			-19,
			-11,
			10,
			-1,
			-1,
		},
		// left > right
		{
			10,
			5,
			10,
			-1,
			-1,
		},
		{
			-2,
			-3,
			10,
			-1,
			-1,
		},
		{
			-1,
			0,
			10,
			-1,
			-1,
		},
		{
			-10,
			0,
			10,
			0,
			0,
		},
	}

	for idx, testCase := range testCases {
		actualLeft, actualRight := NewSlice(testCase.indexLeft, testCase.indexRight).Resolve(testCase.size)
		if !assert.Equal(t, testCase.expectedLeft, actualLeft) || !assert.Equal(t, testCase.expectedRight, actualRight) {
			t.Fatalf("ResolveRaw result differ. idx=%d, expected=(%d, %d), actual=(%d, %d)", idx, testCase.expectedLeft, testCase.expectedRight, actualLeft, actualRight)
		}
	}
}
