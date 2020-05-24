package command

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lxdlam/vertex/pkg/container"

	"github.com/lxdlam/vertex/pkg/protocol"
)

type CommandType int

const (
	_ CommandType = iota
	SystemCommandType
	AccessCommandType
	ModifyCommandType
	CommandTypeLength
)

var (
	// ErrCommandNotExist will be raised if the command is not exist
	ErrCommandNotExist = errors.New("command: command not exist")

	// ErrArgumentInvalid will be raised if the Validate call by the command
	ErrArgumentInvalid = errors.New("command: the argument is invalid")
)

type Command interface {
	Name() string
	ParseArguments([]protocol.RedisObject) error

	Execute()
	Result() (protocol.RedisObject, error)

	Cluster() int

	Keys() []string
	ShouldCreate() bool
	SetAccessObjects([]container.ContainerObject)
	TargetContainerType() container.ContainerType

	ToLog() string
	Type() CommandType
}

var keyMap map[string]func(string, int, []protocol.RedisObject) Command = nil
var lock sync.RWMutex

func init() {
	lock.Lock()
	defer lock.Unlock()

	keyMap = make(map[string]func(string, int, []protocol.RedisObject) Command)
	keyMap["SET"] = newGlobalCommand
	keyMap["GET"] = newGlobalCommand
}

// NewCommand will returns a new command by the name
func NewCommand(name string, index int, arguments []protocol.RedisObject) (Command, error) {
	lock.RLock()
	defer lock.RUnlock()

	fn, ok := keyMap[name]
	if !ok {
		return nil, fmt.Errorf("command name not found")
	}

	return fn(name, index, arguments), nil
}
