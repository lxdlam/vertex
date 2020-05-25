package protocol_test

import (
	"strings"
	"testing"

	. "github.com/lxdlam/vertex/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

// A successful tests are adapted from resp_test.go
func TestSuccessParses(t *testing.T) {
	testCases := []string{
		"$0\r\n\r\n",                  // EmptyBulkString
		"*0\r\n",                      // EmptyArray
		"$-1\r\n",                     // NullBulkString
		"*-1\r\n",                     // NullArray
		"+OK\r\n",                     // SimpleString
		"-Error Message\r\n",          // Error
		":1000\r\n",                   // Integers
		":0\r\n",                      // Zero
		":-10000\r\n",                 // Negative integers
		"$14\r\nHello \rWorld!\n\r\n", // Bulk String
		// Complex BulkString
		"$206\r\n$198\r\nMcpJxQSaoHhgLbAsdPTRjGFCZqrwlEivnIuKkfOVYzWDNXmyUetB\r\n\r\n\r\x0c \n\t\x0b\r\n\r\n^|,!%#`:_<*-=?;$[\\&.>\"(~+]'/}{)@\r\n\r\nETbgCoRcXOWezIHKxmQqVvhADyGZJplufNFjakdYPrwMLsBitUSn\r\n\r\n_)~&!\\.[<=*]^>+;$%(/'@,}\"{:-|#`?\r\n\r\n\r\n \t\r\n\r\n",
		// Array
		"*4\r\n+OK\r\n-Error Message\r\n:1000\r\n$14\r\nHello \rWorld!\n\r\n",
		// NullBulkString in array
		"*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n",
		// Nested array
		"*1\r\n*1\r\n*1\r\n*2\r\n:-1\r\n+OK\r\n",
		// A test case mixed almost everything
		"*14\r\n+OK\r\n-Error Message\r\n-Go Error\r\n:1000\r\n:0\r\n:-10000\r\n$0\r\n\r\n$14\r\nHello \rWorld!\n\r\n$206\r\n$198\r\nMcpJxQSaoHhgLbAsdPTRjGFCZqrwlEivnIuKkfOVYzWDNXmyUetB\r\n\r\n\r\x0c \n\t\x0b\r\n\r\n^|,!%#`:_<*-=?;$[\\&.>\"(~+]'/}{)@\r\n\r\nETbgCoRcXOWezIHKxmQqVvhADyGZJplufNFjakdYPrwMLsBitUSn\r\n\r\n_)~&!\\.[<=*]^>+;$%(/'@,}\"{:-|#`?\r\n\r\n\r\n \t\r\n\r\n*0\r\n*4\r\n+OK\r\n-Error Message\r\n:1000\r\n$14\r\nHello \rWorld!\n\r\n*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n*1\r\n*1\r\n*1\r\n*2\r\n:-1\r\n+OK\r\n*-1\r\n",
	}

	for idx, testCase := range testCases {
		if !testSimpleParse(t, testCase) {
			t.Fatalf("test %d failed. case=%s", idx, testCase)
		}
	}
}

func testSimpleParse(t *testing.T, raw string) bool {
	obj, err := Parse(strings.NewReader(raw))

	return assert.Nil(t, err) && assert.Equal(t, raw, obj.String())
}
