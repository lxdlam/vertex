package command

import "github.com/lxdlam/vertex/pkg/protocol"

func newSetCommand(name string, index int, arguments []protocol.RedisObject) (Command, error) {

	return nil, ErrCommandNotExist
}
