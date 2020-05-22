package util_test

import (
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
