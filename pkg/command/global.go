package command

import (
	"fmt"

	"github.com/lxdlam/vertex/pkg/container"
	"github.com/lxdlam/vertex/pkg/protocol"
)

func newGlobalCommand(name string, index int, arguments []protocol.RedisObject) Command {
	switch name {
	case "SET":
		s := &setCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		if err == nil {
			return s
		}

		return nil
	case "GET":
		g := &getCommand{
			index: index,
		}
		err := g.ParseArguments(arguments)
		if err == nil {
			return g
		}

		return nil
	}

	return nil
}

type setCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	arguments    protocol.RedisString
	result       protocol.RedisString
	err          error
}

func (s *setCommand) Name() string {
	return "SET"
}

func (s *setCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 2 {
		return fmt.Errorf("invalid argument size, expected 2. got=%d", len(objects))
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return fmt.Errorf("parse key failed. raw=%+v", objects[0])
	}

	s.key = tmpObj.Data()

	s.arguments, ok = objects[1].(protocol.RedisString)
	if !ok {
		return fmt.Errorf("parse argument failed. raw=%+v", objects[1])
	}

	return nil
}

func (s *setCommand) Execute() {
	if s.accessObject == nil {
		s.err = fmt.Errorf("nil access object")
		return
	}

	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	s.err = s.accessObject.(container.StringMap).Set([]*container.StringContainer{container.NewString(s.key)}, []*container.StringContainer{container.NewString(s.arguments.Data())})

	if s.err == nil {
		s.result = protocol.NewSimpleRedisString("OK")
	}
}

func (s *setCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *setCommand) Cluster() int {
	return s.index
}

func (s *setCommand) ToLog() string {
	panic("implement me")
}

func (s *setCommand) Type() CommandType {
	return ModifyCommandType
}

func (s *setCommand) Keys() []string {
	return []string{s.key}
}

func (s *setCommand) ShouldCreate() bool {
	return true
}

func (s *setCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *setCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}

type getCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisString
	err          error
}

func (g *getCommand) Name() string {
	return "GET"
}

func (g *getCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return fmt.Errorf("invalid argument size, expected 2. got=%d", len(objects))
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return fmt.Errorf("parse key failed. raw=%+v", objects[0])
	}

	g.key = tmpObj.Data()

	return nil
}

func (g *getCommand) Execute() {
	if g.accessObject == nil {
		g.err = fmt.Errorf("nil access object")
		return
	}

	if g.accessObject.Type() != g.TargetContainerType() {
		g.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", g.TargetContainerType(), g.accessObject.Type())
		return
	}

	ret := g.accessObject.(container.StringMap).Get([]*container.StringContainer{container.NewString(g.key)})

	if len(ret) == 1 && ret[0] != nil {
		g.result = protocol.NewBulkRedisString(ret[0].String())
	} else {
		g.result = protocol.NewNullBulkRedisString()
	}
}

func (g *getCommand) Result() (protocol.RedisObject, error) {
	return g.result, g.err
}

func (g *getCommand) Cluster() int {
	return g.index
}

func (g *getCommand) ToLog() string {
	panic("implement me")
}

func (g *getCommand) Type() CommandType {
	return AccessCommandType
}

func (g *getCommand) Keys() []string {
	return []string{g.key}
}

func (g *getCommand) ShouldCreate() bool {
	return false
}

func (g *getCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	g.accessObject = objects[0]
}

func (g *getCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}
