package internal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/lxdlam/vertex/pkg/protocol"
)

func FormatInput(input string) (string, error) {
	var arr []protocol.RedisObject
	reader := bufio.NewReader(strings.NewReader(input))

	for {
		s, err := parseString(reader)

		if err != nil {
			if errors.Is(err, io.EOF) {
				if s != "" {
					arr = append(arr, protocol.NewBulkRedisString(s))
				}
				break
			} else {
				return "", err
			}
		}

		if s != "" {
			arr = append(arr, protocol.NewBulkRedisString(s))
		}
	}

	return protocol.NewRedisArray(arr).String(), nil
}

func parseString(reader *bufio.Reader) (string, error) {
	var buf bytes.Buffer
	for {
		b, e := reader.ReadByte()
		if e != nil {
			if errors.Is(e, io.EOF) {
				return buf.String(), e
			}
			return "", e
		}

		if b == '"' {
			return parseQuote(reader)
		} else if b == ' ' {
			break
		}

		buf.WriteByte(b)
	}

	return buf.String(), nil
}

func parseQuote(reader *bufio.Reader) (string, error) {
	var buf bytes.Buffer
	for {
		b, e := reader.ReadByte()
		if e != nil {
			return "", e
		} else if b == '"' {
			break
		}

		buf.WriteByte(b)
	}

	return buf.String(), nil
}

func FormatOutput(obj protocol.RedisObject) string {
	if obj == nil {
		panic("nil object")
	}

	switch obj.Type() {
	case protocol.SimpleStringType, protocol.BulkStringType:
		return formatString(obj.(protocol.RedisString))
	case protocol.ErrorType:
		return formatError(obj.(protocol.RedisError))
	case protocol.IntegerType:
		return formatInteger(obj.(protocol.RedisInteger))
	case protocol.ArrayType:
		return formatArray(obj.(protocol.RedisArray))
	default:
		return fmt.Sprintf("ClientError: unknown obj! object=%+v", obj)
	}
}

func formatString(obj protocol.RedisString) string {
	if obj.String() == protocol.NullBulkStringLiteral {
		return "(nil)"
	}

	return fmt.Sprintf("\"%s\"", obj.Data())
}

func formatError(obj protocol.RedisError) string {
	return fmt.Sprintf("Server error: %s", obj.Message())
}

func formatInteger(obj protocol.RedisInteger) string {
	return fmt.Sprintf("(integer) %d", obj.Data())
}

func formatArray(obj protocol.RedisArray) string {
	if obj.String() == protocol.NullArrayLiteral {
		return "(empty list or set)"
	}

	var ret []string

	for idx, item := range obj.Data() {
		ret = append(ret, fmt.Sprintf("%d) %s", idx+1, FormatOutput(item)))
	}

	return strings.Join(ret, "\n")
}
