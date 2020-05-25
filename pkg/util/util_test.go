package util_test

import (
	"errors"
	"testing"

	. "github.com/lxdlam/vertex/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestLexicalCompare(t *testing.T) {
	testCases := []struct {
		lhs, rhs string
		expected int
	}{
		{
			"Zebra",
			"ant",
			-1,
		},
		{
			"Apple",
			"apple",
			-1,
		},
		{
			"orange",
			"apple",
			1,
		},
		{
			"apple",
			"apple",
			0,
		},
		{
			"maple",
			"morning",
			-2,
		},
		{
			"apple",
			"apple",
			0,
		},
		{
			"applecart",
			"apple",
			5,
		},
		{
			"face",
			"facebook",
			-4,
		},
	}

	for _, testCase := range testCases {
		actual := LexicalCompare(testCase.lhs, testCase.rhs)
		assert.Equal(t, testCase.expected, actual)
	}
}

func TestParseInt64(t *testing.T) {
	testCases := []struct {
		str           string
		expected      int64
		expectedError error
	}{
		{
			"1",
			1,
			nil,
		},
		{
			"0",
			0,
			nil,
		},
		{
			"-123456789",
			-123456789,
			nil,
		},
		{
			"123123123123123123123",
			0,
			errors.New("should error"),
		},
		{
			"",
			0,
			errors.New("should error"),
		},
		{
			"-",
			0,
			errors.New("should error"),
		},
		{
			"+",
			0,
			errors.New("should error"),
		},
		{
			"+123456",
			123456,
			nil,
		},
		{
			"+0",
			0,
			nil,
		},
		{
			"-0",
			0,
			nil,
		},
		{
			"1a",
			0,
			errors.New("should error"),
		},
	}

	for _, testCase := range testCases {
		ret, err := ParseInt64(testCase.str)
		if testCase.expectedError != nil {
			assert.NotNil(t, err)
		} else {
			assert.Equal(t, testCase.expected, ret)
		}
	}
}
