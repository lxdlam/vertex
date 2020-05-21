package replication

import (
	"github.com/gorilla/websocket"
	"github.com/lxdlam/vertex/pkg/util"
	"net"
)

type ReplicationCenter interface {
}

type node struct {
	id        string
	addr      net.Addr
	timeStamp int64
	conn      *websocket.Conn
}

type replicationCenter struct {
	master    *node
	clients   []*node
	localAddr net.Addr
}

func newNode(conn *websocket.Conn) *node {
	return &node{
		id:        util.GenNewUUID(),
		addr:      conn.RemoteAddr(),
		timeStamp: -1, // never syncs
		conn:      conn,
	}
}
