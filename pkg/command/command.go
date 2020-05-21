package command

import (
	"github.com/lxdlam/vertex/pkg/protocol"
)

type Command interface {
	Name() string
	SetArguments([]*protocol.RedisObject)
	Execute(string) (protocol.RedisObject, error)
	ToLog() string
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

type AccessCommand interface {
	OperationCommand
}

var keyMap = map[string]func([]*protocol.RedisObject) Command{
}

type CommandFactory interface {
	NewCommand(string, []*protocol.RedisObject) (Command, error)
}
