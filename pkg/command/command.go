package command

import (
	"github.com/lxdlam/vertex/pkg/container"
	"github.com/lxdlam/vertex/pkg/protocol"
)

type Command interface {
	Name() string
	Execute(container.Container) (protocol.RedisObject, error)
	Validate() bool
}

type OperationCommand interface {
	Command

	Key() string
	Arguments() []*protocol.RedisObject
}

type ModifyCommand interface {
	OperationCommand

	GenCancelCommand() Command
}
