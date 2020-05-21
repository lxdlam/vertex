package command

import (
	"errors"
	"github.com/lxdlam/vertex/pkg/protocol"
)

var (
	// ErrCommandNotExist will be raised if the command is not exist
	ErrCommandNotExist = errors.New("command: command not exist")

	// ErrArgumentInvalid will be raised if the Validate call by the command
	ErrArgumentInvalid = errors.New("command: the argument is invalid")
)

var keyMap map[string]func(int, []*protocol.RedisObject) Command

type Command interface {
	Name() string
	ParseArguments(int, []*protocol.RedisObject) (bool, error)

	ToLog() string
}

// NewCommand
func NewCommand(key string, arguments []*protocol.RedisObject) (Command, error) {

}
