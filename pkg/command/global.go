package command

import (
	"errors"
	"fmt"

	"github.com/lxdlam/vertex/pkg/util"

	"github.com/lxdlam/vertex/pkg/container"
	"github.com/lxdlam/vertex/pkg/protocol"
)

func newGlobalCommand(name string, index int, arguments []protocol.RedisObject) (Command, error) {
	switch name {
	case "set":
		s := &setCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "get":
		g := &getCommand{
			index: index,
		}
		err := g.ParseArguments(arguments)
		return g, err
	case "mset":
		s := &msetCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "mget":
		g := &mgetCommand{
			index: index,
		}
		err := g.ParseArguments(arguments)
		return g, err
	case "exists":
		e := &existsCommand{
			index: index,
		}
		err := e.ParseArguments(arguments)
		return e, err
	case "strlen":
		s := &strlenCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "append":
		a := &appendCommand{
			index: index,
		}
		err := a.ParseArguments(arguments)
		return a, err
	case "incr":
		var err error
		i := &incrCommand{
			index: index,
		}
		if len(arguments) >= 1 {
			err = i.ParseArguments(arguments[:1])
		} else {
			err = i.ParseArguments(arguments)
		}
		return i, err
	case "incrby":
		i := &incrCommand{
			index: index,
		}
		err := i.ParseArguments(arguments)
		return i, err
	case "decr":
		var err error
		i := &decrCommand{
			index: index,
		}
		if len(arguments) >= 1 {
			err = i.ParseArguments(arguments[:1])
		} else {
			err = i.ParseArguments(arguments)
		}
		return i, err
	case "decrby":
		i := &decrCommand{
			index: index,
		}
		err := i.ParseArguments(arguments)
		return i, err
	case "getrange":
		g := &getRangeCommand{
			index: index,
		}
		err := g.ParseArguments(arguments)
		return g, err
	}

	return nil, ErrCommandNotExist
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
	return "set"
}

func (s *setCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	s.key = tmpObj.Data()

	s.arguments, ok = objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
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

	if s.arguments == nil {
		s.err = ErrArgumentInvalid
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
	return "get"
}

func (g *getCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
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

type msetCommand struct {
	index        int
	accessObject container.ContainerObject
	keys         []string
	values       []string
	result       protocol.RedisObject
	err          error
}

func (s *msetCommand) Name() string {
	return "mset"
}

func (s *msetCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)

	if length%2 != 0 {
		return ErrArgumentInvalid
	}

	for idx := 0; idx < length; idx += 2 {
		keyObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		s.keys = append(s.keys, keyObj.Data())

		valueObj, ok := objects[idx+1].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		s.values = append(s.values, valueObj.Data())
	}

	return nil
}

func (s *msetCommand) Execute() {
	if s.accessObject == nil {
		s.err = fmt.Errorf("nil access object")
		return
	}

	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	if s.keys == nil || s.values == nil {
		s.err = ErrArgumentInvalid
		return
	}

	var keys, values []*container.StringContainer
	length := len(s.keys)

	for idx := 0; idx < length; idx++ {
		keys = append(keys, container.NewString(s.keys[idx]))
		values = append(values, container.NewString(s.values[idx]))
	}

	// MSET will produce no error.
	_ = s.accessObject.(container.StringMap).Set(keys, values)

	s.err = nil
	s.result = protocol.NewSimpleRedisString("OK")
}

func (s *msetCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *msetCommand) Cluster() int {
	return s.index
}

func (s *msetCommand) ToLog() string {
	panic("implement me")
}

func (s *msetCommand) Type() CommandType {
	return ModifyCommandType
}

func (s *msetCommand) Keys() []string {
	return s.keys
}

func (s *msetCommand) ShouldCreate() bool {
	return true
}

func (s *msetCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *msetCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}

type mgetCommand struct {
	keys         []string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisArray
	err          error
}

func (g *mgetCommand) Name() string {
	return "mget"
}

func (g *mgetCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)

	if length <= 0 {
		return ErrArgumentInvalid
	}

	for idx := 0; idx < length; idx++ {
		tmpObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		g.keys = append(g.keys, tmpObj.Data())
	}

	return nil
}

func (g *mgetCommand) Execute() {
	if g.accessObject == nil {
		g.err = fmt.Errorf("nil access object")
		return
	}

	if g.accessObject.Type() != g.TargetContainerType() {
		g.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", g.TargetContainerType(), g.accessObject.Type())
		return
	}

	var keys []*container.StringContainer

	for _, key := range g.keys {
		keys = append(keys, container.NewString(key))
	}

	var values []protocol.RedisObject

	ret := g.accessObject.(container.StringMap).Get(keys)
	length := len(ret)

	for idx := 0; idx < length; idx++ {
		if ret[idx] == nil {
			values = append(values, protocol.NewNullBulkRedisString())
		} else {
			values = append(values, protocol.NewBulkRedisString(ret[idx].String()))
		}
	}

	g.result = protocol.NewRedisArray(values)
}

func (g *mgetCommand) Result() (protocol.RedisObject, error) {
	return g.result, g.err
}

func (g *mgetCommand) Cluster() int {
	return g.index
}

func (g *mgetCommand) ToLog() string {
	panic("implement me")
}

func (g *mgetCommand) Type() CommandType {
	return AccessCommandType
}

func (g *mgetCommand) Keys() []string {
	return g.keys
}

func (g *mgetCommand) ShouldCreate() bool {
	return false
}

func (g *mgetCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	g.accessObject = objects[0]
}

func (g *mgetCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}

type existsCommand struct {
	keys         []string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (e *existsCommand) Name() string {
	return "exists"
}

func (e *existsCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)

	if length <= 0 {
		return ErrArgumentInvalid
	}

	for idx := 0; idx < length; idx++ {
		tmpObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		e.keys = append(e.keys, tmpObj.Data())
	}

	return nil
}

func (e *existsCommand) Execute() {
	if e.accessObject == nil {
		e.err = fmt.Errorf("nil access object")
		return
	}

	if e.accessObject.Type() != e.TargetContainerType() {
		e.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", e.TargetContainerType(), e.accessObject.Type())
		return
	}

	var keys []*container.StringContainer

	for _, key := range e.keys {
		keys = append(keys, container.NewString(key))
	}

	ret := e.accessObject.(container.StringMap).Exists(keys)

	e.result = protocol.NewRedisInteger(int64(ret))
}

func (e *existsCommand) Result() (protocol.RedisObject, error) {
	return e.result, e.err
}

func (e *existsCommand) Cluster() int {
	return e.index
}

func (e *existsCommand) ToLog() string {
	panic("implement me")
}

func (e *existsCommand) Type() CommandType {
	return AccessCommandType
}

func (e *existsCommand) Keys() []string {
	return e.keys
}

func (e *existsCommand) ShouldCreate() bool {
	return false
}

func (e *existsCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	e.accessObject = objects[0]
}

func (e *existsCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}

type strlenCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (s *strlenCommand) Name() string {
	return "strlen"
}

func (s *strlenCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	s.key = tmpObj.Data()

	return nil
}

func (s *strlenCommand) Execute() {
	if s.accessObject == nil {
		s.err = fmt.Errorf("nil access object")
		return
	}

	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	ret, _ := s.accessObject.(container.StringMap).StringLen(container.NewString(s.key))

	s.result = protocol.NewRedisInteger(int64(ret))
}

func (s *strlenCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *strlenCommand) Cluster() int {
	return s.index
}

func (s *strlenCommand) Keys() []string {
	return []string{s.key}
}

func (s *strlenCommand) ShouldCreate() bool {
	return false
}

func (s *strlenCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *strlenCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}

func (s *strlenCommand) ToLog() string {
	panic("implement me")
}

func (s *strlenCommand) Type() CommandType {
	return AccessCommandType
}

type appendCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	arguments    protocol.RedisString
	result       protocol.RedisInteger
	err          error
}

func (a *appendCommand) Name() string {
	return "append"
}

func (a *appendCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	a.key = tmpObj.Data()

	a.arguments, ok = objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	return nil
}

func (a *appendCommand) Execute() {
	if a.accessObject == nil {
		a.err = fmt.Errorf("nil access object")
		return
	}

	if a.accessObject.Type() != a.TargetContainerType() {
		a.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", a.TargetContainerType(), a.accessObject.Type())
		return
	}

	if a.arguments == nil {
		a.err = ErrArgumentInvalid
		return
	}

	ret := a.accessObject.(container.StringMap).Append(container.NewString(a.key), container.NewString(a.arguments.Data()))

	a.result = protocol.NewRedisInteger(int64(ret))
	a.err = nil
}

func (a *appendCommand) Result() (protocol.RedisObject, error) {
	return a.result, a.err
}

func (a *appendCommand) Cluster() int {
	return a.index
}

func (a *appendCommand) Keys() []string {
	return []string{a.key}
}

func (a *appendCommand) ShouldCreate() bool {
	return true
}

func (a *appendCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	a.accessObject = objects[0]
}

func (a *appendCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}

func (a *appendCommand) ToLog() string {
	panic("implement me")
}

func (a *appendCommand) Type() CommandType {
	return ModifyCommandType
}

type incrCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	arguments    protocol.RedisString
	result       protocol.RedisInteger
	err          error
}

func (i *incrCommand) Name() string {
	return "incrby"
}

func (i *incrCommand) ParseArguments(objects []protocol.RedisObject) error {
	l := len(objects)
	if l == 0 || l > 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	i.key = tmpObj.Data()

	if l == 1 {
		i.arguments = protocol.NewSimpleRedisString("1")
	} else {
		i.arguments, ok = objects[1].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}
	}

	return nil
}

func (i *incrCommand) Execute() {
	if i.accessObject == nil {
		i.err = fmt.Errorf("nil access object")
		return
	}

	if i.accessObject.Type() != i.TargetContainerType() {
		i.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", i.TargetContainerType(), i.accessObject.Type())
		return
	}

	if i.arguments == nil {
		i.err = ErrArgumentInvalid
		return
	}

	var ret int64

	increment, err := util.ParseInt64(i.arguments.Data())
	if err == nil {
		ret, i.err = i.accessObject.(container.StringMap).Increase(container.NewString(i.key), increment)
	} else {
		i.err = err
	}

	if i.err == nil {
		i.result = protocol.NewRedisInteger(ret)
	} else {
		if errors.Is(i.err, container.ErrKeyNotFound) {
			i.err = ErrNoSuchKey
		}
	}
}

func (i *incrCommand) Result() (protocol.RedisObject, error) {
	return i.result, i.err
}

func (i *incrCommand) Cluster() int {
	return i.index
}

func (i *incrCommand) ToLog() string {
	panic("implement me")
}

func (i *incrCommand) Type() CommandType {
	return ModifyCommandType
}

func (i *incrCommand) Keys() []string {
	return []string{i.key}
}

func (i *incrCommand) ShouldCreate() bool {
	return true
}

func (i *incrCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	i.accessObject = objects[0]
}

func (i *incrCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}

type decrCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	arguments    protocol.RedisString
	result       protocol.RedisInteger
	err          error
}

func (d *decrCommand) Name() string {
	return "decrby"
}

func (d *decrCommand) ParseArguments(objects []protocol.RedisObject) error {
	l := len(objects)
	if l == 0 || l > 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	d.key = tmpObj.Data()

	if l == 1 {
		d.arguments = protocol.NewSimpleRedisString("1")
	} else {
		d.arguments, ok = objects[1].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}
	}

	return nil
}

func (d *decrCommand) Execute() {
	if d.accessObject == nil {
		d.err = fmt.Errorf("nil access object")
		return
	}

	if d.accessObject.Type() != d.TargetContainerType() {
		d.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", d.TargetContainerType(), d.accessObject.Type())
		return
	}

	if d.arguments == nil {
		d.err = ErrArgumentInvalid
		return
	}

	var ret int64

	decrement, err := util.ParseInt64(d.arguments.Data())
	if err == nil {
		ret, d.err = d.accessObject.(container.StringMap).Decrease(container.NewString(d.key), decrement)
	} else {
		d.err = err
	}

	if d.err == nil {
		d.result = protocol.NewRedisInteger(ret)
	} else {
		if errors.Is(d.err, container.ErrKeyNotFound) {
			d.err = ErrNoSuchKey
		}
	}
}

func (d *decrCommand) Result() (protocol.RedisObject, error) {
	return d.result, d.err
}

func (d *decrCommand) Cluster() int {
	return d.index
}

func (d *decrCommand) ToLog() string {
	panic("implement me")
}

func (d *decrCommand) Type() CommandType {
	return ModifyCommandType
}

func (d *decrCommand) Keys() []string {
	return []string{d.key}
}

func (d *decrCommand) ShouldCreate() bool {
	return true
}

func (d *decrCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	d.accessObject = objects[0]
}

func (d *decrCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}

type getRangeCommand struct {
	key          string
	index        int
	start        int
	end          int
	accessObject container.ContainerObject
	result       protocol.RedisString
	err          error
}

func (g *getRangeCommand) Name() string {
	return "getrange"
}

func (g *getRangeCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 3 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	g.key = tmpObj.Data()

	startObj, ok := objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	tmpIdx, err := util.ParseInt64(startObj.Data())
	if err != nil {
		return ErrArgumentInvalid
	}

	g.start = int(tmpIdx)

	endObj, ok := objects[2].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	tmpIdx, err = util.ParseInt64(endObj.Data())
	if err != nil {
		return ErrArgumentInvalid
	}

	g.end = int(tmpIdx)

	return nil
}

func (g *getRangeCommand) Execute() {
	if g.accessObject == nil {
		g.err = fmt.Errorf("nil access object")
		return
	}

	if g.accessObject.Type() != g.TargetContainerType() {
		g.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", g.TargetContainerType(), g.accessObject.Type())
		return
	}

	ret, err := g.accessObject.(container.StringMap).GetRange(container.NewString(g.key), g.start, g.end)

	if err != nil {
		g.result = protocol.NewBulkRedisString("")
	} else {
		g.result = protocol.NewBulkRedisString(ret.String())
	}
}

func (g *getRangeCommand) Result() (protocol.RedisObject, error) {
	return g.result, g.err
}

func (g *getRangeCommand) Cluster() int {
	return g.index
}

func (g *getRangeCommand) ToLog() string {
	panic("implement me")
}

func (g *getRangeCommand) Type() CommandType {
	return AccessCommandType
}

func (g *getRangeCommand) Keys() []string {
	return []string{g.key}
}

func (g *getRangeCommand) ShouldCreate() bool {
	return false
}

func (g *getRangeCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	g.accessObject = objects[0]
}

func (g *getRangeCommand) TargetContainerType() container.ContainerType {
	return container.GlobalType
}
