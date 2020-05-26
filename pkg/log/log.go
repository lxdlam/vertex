package log

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lxdlam/vertex/pkg/protocol"
	"github.com/lxdlam/vertex/pkg/util"
	"io"
	"strconv"
	"time"
)

var (
	// ErrEmptyLog will be raised if in PackLog the proto.Marshal returns a zero length byte slice
	ErrEmptyLog = errors.New("empty vertex log message")
)

var (
	errUnexpectedEOF  = errors.New("unexpected eof")
	errInvalidMessage = errors.New("invalid proto message")
)

// NewLog will create a new VertexLog instance
func NewLog(name string, index int, objects []protocol.RedisObject) *VertexLog {
	vl := &VertexLog{}

	vl.Id = util.GenNewUUID()
	vl.Time = time.Now().UTC().UnixNano()
	vl.Host = util.GetIP()
	vl.Name = name
	vl.Index = int32(index)

	var arguments []string
	for _, item := range objects[1:] {
		arguments = append(arguments, item.String())
	}

	vl.Arguments = arguments
	vl.RawRequest = protocol.NewRedisArray(objects).String()

	return vl
}

func FormatLog(vl *VertexLog) string {
	return fmt.Sprintf("VectexLog{id=%s, time=%d, host=%s, name=%s, index=%d, arguments=[%s], raw_request=%s}", vl.Id, vl.Time, vl.Host, vl.Name, vl.Index, util.QuoteJoin(vl.Arguments, ","), strconv.Quote(vl.RawRequest))
}

// PackLog will package a log into the below form:
//     represent: | length | log |
//     bytes:         ^4      ^the value in length section
func PackLog(vl *VertexLog) ([]byte, error) {
	var buf bytes.Buffer

	message, err := proto.Marshal(vl)
	if err != nil {
		return nil, err
	}

	length := len(message)

	if length == 0 {
		return nil, ErrEmptyLog
	}

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(length))

	buf.Write(b)
	buf.Write(message)

	return buf.Bytes(), nil
}

type logReader struct {
	reader *bufio.Reader
}

// ParseLog will parse a log from a io.Reader which can be a file or net stream
func ParseLog(reader io.Reader) ([]*VertexLog, error) {
	var ret []*VertexLog

	lr := logReader{
		reader: bufio.NewReader(reader),
	}

	for {
		obj, err := lr.readLog()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else if errors.Is(err, errUnexpectedEOF) {
				return ret, fmt.Errorf("unexpected eof when parsing log")
			} else if errors.Is(err, errInvalidMessage) {
				// Try parse more
				continue
			} else {
				return nil, fmt.Errorf("unexpected error met. err={%w}", err)
			}
		}

		ret = append(ret, obj)
	}

	return ret, nil
}

func (lr *logReader) readBytes(count int) ([]byte, error) {
	var ret []byte

	for idx := 0; idx < count; idx++ {
		b, err := lr.reader.ReadByte()
		if err != nil {
			return nil, err
		}

		ret = append(ret, b)
	}

	return ret, nil
}

func (lr *logReader) readLog() (*VertexLog, error) {
	lengthByte, err := lr.readBytes(4)
	if errors.Is(err, io.EOF) {
		if len(lengthByte) != 0 {
			return nil, errUnexpectedEOF
		}
		return nil, io.EOF
	} else if err != nil {
		return nil, err
	}

	length := int(binary.LittleEndian.Uint32(lengthByte))

	logByte, err := lr.readBytes(length)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errUnexpectedEOF
		}

		return nil, err
	}

	v := &VertexLog{}

	if err := proto.Unmarshal(logByte, v); err != nil {
		return nil, errInvalidMessage
	}

	return v, nil
}
