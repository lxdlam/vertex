package network

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/lxdlam/vertex/pkg/protocol"
	"github.com/lxdlam/vertex/pkg/util"
)

var (
	defaultExpireTime = 30 * time.Hour

	// ErrConnIsClosed will be raised if do any operation on a closed conn
	ErrConnIsClosed = errors.New("conn: conn is already closed")
)

// Conn is a client connection object that will handle a expire time
type Conn interface {
	Read() (protocol.RedisObject, error)
	Write(string) error
	Close() error

	IsClosed() bool

	Addr() string
	ID() string
}

type conn struct {
	id         string
	addr       string
	tcpConn    net.Conn
	expireTime time.Duration
	closed     int32
	resetChan  chan byte
	closeChan  chan struct{}
}

// NewConn equals NewConnWithExpire(conn, defaultExpireTime)
func NewConn(conn net.Conn) Conn {
	return NewConnWithExpire(conn, defaultExpireTime)
}

// NewConnWithExpire will return a new Conn with the expire time set to expireTime. The
// expire time is global, what means if no operation happens from the last operation for
// expire time, it will close the operation and do not serve anymore.
func NewConnWithExpire(tcpConn net.Conn, expireTime time.Duration) Conn {
	c := &conn{
		id:         util.GenNewUUID(),
		addr:       tcpConn.RemoteAddr().String(),
		tcpConn:    tcpConn,
		expireTime: expireTime,
		closed:     0,
		resetChan:  make(chan byte),
		closeChan:  make(chan struct{}),
	}

	c.startExpireWorker()

	return c
}

func (c *conn) Read() (protocol.RedisObject, error) {
	c.resetChan <- 0

	var buf bytes.Buffer
	// 1M buffer
	readBuf := make([]byte, 1024*1024)

	_, err := c.tcpConn.Read(readBuf)
	if errors.Is(err, io.EOF) {
		return nil, ErrConnIsClosed
	} else if err != nil {
		return nil, fmt.Errorf("conn: read with unexpected error. conn.id=%s, conn.tcpConn.addr=%s, err={%w}", c.id, c.addr, err)
	}

	buf.Write(readBuf)

	obj, err := protocol.Parse(&buf)

	if err != nil {
		return nil, fmt.Errorf("conn: invalid resp input. conn.id=%s, conn.tcpConn.addr=%s, err={%w}", c.id, c.addr, err)
	}

	return obj, nil
}

func (c *conn) Write(s string) error {
	if c.IsClosed() {
		return ErrConnIsClosed
	}

	c.resetChan <- 0

	n, err := c.tcpConn.Write([]byte(s))

	if err != nil {
		return fmt.Errorf("conn: write a response met an error. s=%s, conn.id=%s, conn.tcpConn.addr=%s, err={%w}", s, c.id, c.addr, err)
	} else if n != len(s) {
		return fmt.Errorf("conn: not all bytes writes to the tcpConn. s=%s, len(s)=%d, n=%d, conn.id=%s, conn.tcpConn.addr=%s, err={%w}", s, len(s), n, c.id, c.addr, err)
	}

	return nil
}

func (c *conn) Close() error {
	if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		// discard all streams
		_ = c.tcpConn.Close()
		close(c.closeChan)
		return nil
	}

	return ErrConnIsClosed
}

func (c *conn) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

func (c *conn) Addr() string {
	return c.addr
}

func (c *conn) ID() string {
	return c.id
}

func (c *conn) startExpireWorker() {
	go func() {
	Outer:
		for {
			select {
			case <-c.resetChan:
				break
			case <-c.closeChan:
				break Outer
			case <-time.After(c.expireTime):
				_ = c.Close()
				break Outer
			}
		}
	}()
}
