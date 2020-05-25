package command

import (
	"errors"
	"strings"
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

	// ErrNoSuchKey will be raised when access an non-exist container
	ErrNoSuchKey = errors.New("command: no such key")
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

var keyMap map[string]func(string, int, []protocol.RedisObject) (Command, error) = nil
var lock sync.RWMutex

func init() {
	lock.Lock()
	defer lock.Unlock()

	keyMap = make(map[string]func(string, int, []protocol.RedisObject) (Command, error))

	// Global Commands
	keyMap["set"] = newGlobalCommand
	keyMap["get"] = newGlobalCommand
	keyMap["mset"] = newGlobalCommand
	keyMap["mget"] = newGlobalCommand
	keyMap["exists"] = newGlobalCommand
	keyMap["strlen"] = newGlobalCommand
	keyMap["append"] = newGlobalCommand
	keyMap["incr"] = newGlobalCommand
	keyMap["incrby"] = newGlobalCommand
	keyMap["decr"] = newGlobalCommand
	keyMap["decrby"] = newGlobalCommand
	keyMap["getrange"] = newGlobalCommand

	// List Commands
	keyMap["lpop"] = newListCommand
	keyMap["rpop"] = newListCommand
	keyMap["lindex"] = newListCommand
	keyMap["linsert"] = newListCommand
	keyMap["llen"] = newListCommand
	keyMap["lpush"] = newListCommand
	keyMap["rpush"] = newListCommand
	keyMap["lrange"] = newListCommand
	keyMap["ltrim"] = newListCommand
	keyMap["lrem"] = newListCommand
	keyMap["lset"] = newListCommand

	// Hash Commands
	keyMap["hset"] = newHashCommand
	keyMap["hget"] = newHashCommand
	keyMap["hexists"] = newHashCommand
	keyMap["hdel"] = newHashCommand
	keyMap["hmget"] = newHashCommand
	keyMap["hkeys"] = newHashCommand
	keyMap["hvals"] = newHashCommand
	keyMap["hgetall"] = newHashCommand
	keyMap["hstrlen"] = newHashCommand
	keyMap["hlen"] = newHashCommand

	// Set Commands
	keyMap["sadd"] = newSetCommand
	keyMap["srem"] = newSetCommand
	keyMap["sismember"] = newSetCommand
	keyMap["smembers"] = newSetCommand
	keyMap["srandmember"] = newSetCommand
	keyMap["spop"] = newSetCommand
	keyMap["sdiff"] = newSetCommand
	keyMap["sinter"] = newSetCommand
	keyMap["sunion"] = newSetCommand
	keyMap["scard"] = newSetCommand
}

// NewCommand will returns a new command by the name
func NewCommand(name string, index int, arguments []protocol.RedisObject) (Command, error) {
	lock.RLock()
	defer lock.RUnlock()

	// normalize to lower command name
	name = strings.ToLower(name)

	fn, ok := keyMap[name]
	if !ok {
		return nil, ErrCommandNotExist
	}

	return fn(name, index, arguments)
}
