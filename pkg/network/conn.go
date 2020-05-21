package network

import "net"

type VertexConn struct {
	id      string
	rawConn net.Conn
}
