package replication

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"net"

	"github.com/lxdlam/vertex/pkg/log"
	"github.com/lxdlam/vertex/pkg/util"
)

var ErrIncompleteMessage = errors.New("replica: incomplete message")

// The message layout is like:
//     layout: | length | log data |
//       byte: |    4   |    length's data |
// If the send is failed, length section will be an invalid data.

// Replica will works like a one shot query to the master node.
// It will be redesign later.
type Replica interface {
	Receive() ([]*log.VertexLog, error)
}

type replica struct {
	conn net.Conn
}

func NewReplica(conn net.Conn) Replica {
	if conn == nil {
		return nil
	}

	return &replica{conn: conn}
}

func (r *replica) Receive() ([]*log.VertexLog, error) {
	reader := bufio.NewReader(r.conn)

	lengthBytes, err := util.ReadExactBytes(reader, 4)
	if len(lengthBytes) != 4 {
		return nil, ErrIncompleteMessage
	} else if err != nil {
		return nil, err
	}

	length := int(binary.LittleEndian.Uint32(lengthBytes))

	if length < 0 {
		return nil, ErrIncompleteMessage
	} else if length == 0 {
		return nil, nil
	}

	message, err := util.ReadExactBytes(reader, length)
	if len(message) != length {
		return nil, ErrIncompleteMessage
	} else if err != nil {
		return nil, err
	}

	return log.ParseLog(bytes.NewReader(message))
}
