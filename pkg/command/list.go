package command

import (
	"fmt"
	"github.com/lxdlam/vertex/pkg/container"
	"github.com/lxdlam/vertex/pkg/protocol"
	"github.com/lxdlam/vertex/pkg/util"
	"strings"
)

func newListCommand(name string, index int, arguments []protocol.RedisObject) (Command, error) {
	switch name {
	case "lpop":
		l := &lpopCommand{
			index: index,
		}
		err := l.ParseArguments(arguments)
		return l, err
	case "rpop":
		r := &rpopCommand{
			index: index,
		}
		err := r.ParseArguments(arguments)
		return r, err
	case "lpush":
		l := &lpushCommand{
			index: index,
		}
		err := l.ParseArguments(arguments)
		return l, err
	case "rpush":
		r := &rpushCommand{
			index: index,
		}
		err := r.ParseArguments(arguments)
		return r, err
	case "lrange":
		l := &lrangeCommand{
			index: index,
		}
		err := l.ParseArguments(arguments)
		return l, err
	case "lindex":
		l := &lindexCommand{
			index: index,
		}
		err := l.ParseArguments(arguments)
		return l, err
	case "llen":
		l := &llenCommand{
			index: index,
		}
		err := l.ParseArguments(arguments)
		return l, err
	case "lset":
		l := &lsetCommand{
			index: index,
		}
		err := l.ParseArguments(arguments)
		return l, err
	case "lrem":
		l := &lremCommand{
			index: index,
		}
		err := l.ParseArguments(arguments)
		return l, err
	case "ltrim":
		l := &ltrimCommand{
			index: index,
		}
		err := l.ParseArguments(arguments)
		return l, err
	case "linsert":
		l := &linsertCommand{
			index: index,
		}
		err := l.ParseArguments(arguments)
		return l, err
	}

	return nil, ErrCommandNotExist
}

type lpopCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisString
	err          error
}

func (l *lpopCommand) Name() string {
	return "lpop"
}

func (l *lpopCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.key = tmpObj.Data()

	return nil
}

func (l *lpopCommand) Execute() {
	if l.accessObject == nil {
		l.result = protocol.NewNullBulkRedisString()
		return
	}

	if l.accessObject.Type() != l.TargetContainerType() {
		l.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", l.TargetContainerType(), l.accessObject.Type())
		return
	}

	ret, err := l.accessObject.(container.ListContainer).PopHead()

	if err == nil {
		l.result = protocol.NewBulkRedisString(ret.String())
	} else {
		l.result = protocol.NewNullBulkRedisString()
	}
}

func (l *lpopCommand) Result() (protocol.RedisObject, error) {
	return l.result, l.err
}

func (l *lpopCommand) Cluster() int {
	return l.index
}

func (l *lpopCommand) ToLog() string {
	panic("implement me")
}

func (l *lpopCommand) Type() CommandType {
	return ModifyCommandType
}

func (l *lpopCommand) Keys() []string {
	return []string{l.key}
}

func (l *lpopCommand) ShouldCreate() bool {
	return false
}

func (l *lpopCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	l.accessObject = objects[0]
}

func (l *lpopCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type rpopCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisString
	err          error
}

func (r *rpopCommand) Name() string {
	return "rpop"
}

func (r *rpopCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	r.key = tmpObj.Data()

	return nil
}

func (r *rpopCommand) Execute() {
	if r.accessObject == nil {
		r.result = protocol.NewNullBulkRedisString()
		return
	}

	if r.accessObject.Type() != r.TargetContainerType() {
		r.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", r.TargetContainerType(), r.accessObject.Type())
		return
	}

	ret, err := r.accessObject.(container.ListContainer).PopTail()

	if err == nil {
		r.result = protocol.NewBulkRedisString(ret.String())
	} else {
		r.result = protocol.NewNullBulkRedisString()
	}
}

func (r *rpopCommand) Result() (protocol.RedisObject, error) {
	return r.result, r.err
}

func (r *rpopCommand) Cluster() int {
	return r.index
}

func (r *rpopCommand) ToLog() string {
	panic("implement me")
}

func (r *rpopCommand) Type() CommandType {
	return ModifyCommandType
}

func (r *rpopCommand) Keys() []string {
	return []string{r.key}
}

func (r *rpopCommand) ShouldCreate() bool {
	return false
}

func (r *rpopCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	r.accessObject = objects[0]
}

func (r *rpopCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type lpushCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	arguments    []string
	result       protocol.RedisInteger
	err          error
}

func (l *lpushCommand) Name() string {
	return "lpush"
}

func (l *lpushCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)
	if length <= 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.key = tmpObj.Data()

	for idx := 1; idx < length; idx++ {
		valueObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		l.arguments = append(l.arguments, valueObj.Data())
	}

	return nil
}

func (l *lpushCommand) Execute() {
	if l.accessObject == nil {
		l.err = fmt.Errorf("nil access object")
		return
	}

	if l.accessObject.Type() != l.TargetContainerType() {
		l.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", l.TargetContainerType(), l.accessObject.Type())
		return
	}

	var arguments []*container.StringContainer

	for _, argument := range l.arguments {
		arguments = append(arguments, container.NewString(argument))
	}

	ret, err := l.accessObject.(container.ListContainer).PushHead(arguments)

	l.result = protocol.NewRedisInteger(int64(ret))
	l.err = err
}

func (l *lpushCommand) Result() (protocol.RedisObject, error) {
	return l.result, l.err
}

func (l *lpushCommand) Cluster() int {
	return l.index
}

func (l *lpushCommand) ToLog() string {
	panic("implement me")
}

func (l *lpushCommand) Type() CommandType {
	return ModifyCommandType
}

func (l *lpushCommand) Keys() []string {
	return []string{l.key}
}

func (l *lpushCommand) ShouldCreate() bool {
	return true
}

func (l *lpushCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	l.accessObject = objects[0]
}

func (l *lpushCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type rpushCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	arguments    []string
	result       protocol.RedisInteger
	err          error
}

func (r *rpushCommand) Name() string {
	return "rpush"
}

func (r *rpushCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)
	if length <= 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	r.key = tmpObj.Data()

	for idx := 1; idx < length; idx++ {
		valueObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		r.arguments = append(r.arguments, valueObj.Data())
	}

	return nil
}

func (r *rpushCommand) Execute() {
	if r.accessObject == nil {
		r.err = fmt.Errorf("nil access object")
		return
	}

	if r.accessObject.Type() != r.TargetContainerType() {
		r.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", r.TargetContainerType(), r.accessObject.Type())
		return
	}

	var arguments []*container.StringContainer

	for _, argument := range r.arguments {
		arguments = append(arguments, container.NewString(argument))
	}

	ret, err := r.accessObject.(container.ListContainer).PushTail(arguments)

	r.result = protocol.NewRedisInteger(int64(ret))
	r.err = err
}

func (r *rpushCommand) Result() (protocol.RedisObject, error) {
	return r.result, r.err
}

func (r *rpushCommand) Cluster() int {
	return r.index
}

func (r *rpushCommand) ToLog() string {
	panic("implement me")
}

func (r *rpushCommand) Type() CommandType {
	return ModifyCommandType
}

func (r *rpushCommand) Keys() []string {
	return []string{r.key}
}

func (r *rpushCommand) ShouldCreate() bool {
	return true
}

func (r *rpushCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	r.accessObject = objects[0]
}

func (r *rpushCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type lrangeCommand struct {
	key          string
	index        int
	start        int
	end          int
	accessObject container.ContainerObject
	result       protocol.RedisArray
	err          error
}

func (l *lrangeCommand) Name() string {
	return "lrange"
}

func (l *lrangeCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 3 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.key = tmpObj.Data()

	startObj, ok := objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	tmpIdx, err := util.ParseInt64(startObj.Data())
	if err != nil {
		return ErrArgumentInvalid
	}

	l.start = int(tmpIdx)

	endObj, ok := objects[2].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	tmpIdx, err = util.ParseInt64(endObj.Data())
	if err != nil {
		return ErrArgumentInvalid
	}

	l.end = int(tmpIdx)

	return nil
}

func (l *lrangeCommand) Execute() {
	if l.accessObject == nil {
		l.result = protocol.NewNullRedisArray()
		return
	}

	if l.accessObject.Type() != l.TargetContainerType() {
		l.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", l.TargetContainerType(), l.accessObject.Type())
		return
	}

	ret, err := l.accessObject.(container.ListContainer).Range(l.start, l.end)

	if err != nil {
		l.result = protocol.NewNullRedisArray()
	} else {
		var results []protocol.RedisObject

		for _, item := range ret {
			results = append(results, protocol.NewBulkRedisString(item.String()))
		}

		l.result = protocol.NewRedisArray(results)
	}
}

func (l *lrangeCommand) Result() (protocol.RedisObject, error) {
	return l.result, l.err
}

func (l *lrangeCommand) Cluster() int {
	return l.index
}

func (l *lrangeCommand) ToLog() string {
	panic("implement me")
}

func (l *lrangeCommand) Type() CommandType {
	return AccessCommandType
}

func (l *lrangeCommand) Keys() []string {
	return []string{l.key}
}

func (l *lrangeCommand) ShouldCreate() bool {
	return false
}

func (l *lrangeCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	l.accessObject = objects[0]
}

func (l *lrangeCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type lindexCommand struct {
	key          string
	index        int
	pos          int
	accessObject container.ContainerObject
	result       protocol.RedisString
	err          error
}

func (l *lindexCommand) Name() string {
	return "lindex"
}

func (l *lindexCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.key = tmpObj.Data()

	posObj, ok := objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	tmpIdx, err := util.ParseInt64(posObj.Data())
	if err != nil {
		return ErrArgumentInvalid
	}

	l.pos = int(tmpIdx)

	return nil
}

func (l *lindexCommand) Execute() {
	if l.accessObject == nil {
		l.result = protocol.NewNullBulkRedisString()
		return
	}

	if l.accessObject.Type() != l.TargetContainerType() {
		l.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", l.TargetContainerType(), l.accessObject.Type())
		return
	}

	ret, err := l.accessObject.(container.ListContainer).Index(l.pos)

	if err != nil {
		l.result = protocol.NewNullBulkRedisString()
	} else {
		l.result = protocol.NewBulkRedisString(ret.String())
	}
}

func (l *lindexCommand) Result() (protocol.RedisObject, error) {
	return l.result, l.err
}

func (l *lindexCommand) Cluster() int {
	return l.index
}

func (l *lindexCommand) ToLog() string {
	panic("implement me")
}

func (l *lindexCommand) Type() CommandType {
	return AccessCommandType
}

func (l *lindexCommand) Keys() []string {
	return []string{l.key}
}

func (l *lindexCommand) ShouldCreate() bool {
	return false
}

func (l *lindexCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	l.accessObject = objects[0]
}

func (l *lindexCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type llenCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (l *llenCommand) Name() string {
	return "llen"
}

func (l *llenCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.key = tmpObj.Data()

	return nil
}

func (l *llenCommand) Execute() {
	if l.accessObject == nil {
		l.result = protocol.NewRedisInteger(0)
		return
	}

	if l.accessObject.Type() != l.TargetContainerType() {
		l.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", l.TargetContainerType(), l.accessObject.Type())
		return
	}

	ret := l.accessObject.(container.ListContainer).Len()

	l.result = protocol.NewRedisInteger(int64(ret))
}

func (l *llenCommand) Result() (protocol.RedisObject, error) {
	return l.result, l.err
}

func (l *llenCommand) Cluster() int {
	return l.index
}

func (l *llenCommand) ToLog() string {
	panic("implement me")
}

func (l *llenCommand) Type() CommandType {
	return AccessCommandType
}

func (l *llenCommand) Keys() []string {
	return []string{l.key}
}

func (l *llenCommand) ShouldCreate() bool {
	return false
}

func (l *llenCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	l.accessObject = objects[0]
}

func (l *llenCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type lsetCommand struct {
	key          string
	index        int
	pos          int
	item         protocol.RedisString
	accessObject container.ContainerObject
	result       protocol.RedisString
	err          error
}

func (l *lsetCommand) Name() string {
	return "lset"
}

func (l *lsetCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 3 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.key = tmpObj.Data()

	posObj, ok := objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	tmpIdx, err := util.ParseInt64(posObj.Data())
	if err != nil {
		return ErrArgumentInvalid
	}

	l.pos = int(tmpIdx)

	l.item, ok = objects[2].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	return nil
}

func (l *lsetCommand) Execute() {
	if l.accessObject == nil {
		l.err = ErrNoSuchKey
		return
	}

	if l.accessObject.Type() != l.TargetContainerType() {
		l.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", l.TargetContainerType(), l.accessObject.Type())
		return
	}

	err := l.accessObject.(container.ListContainer).Set(l.pos, container.NewString(l.item.Data()))

	if err != nil {
		l.err = err
	} else {
		l.result = protocol.NewSimpleRedisString("OK")
	}
}

func (l *lsetCommand) Result() (protocol.RedisObject, error) {
	return l.result, l.err
}

func (l *lsetCommand) Cluster() int {
	return l.index
}

func (l *lsetCommand) ToLog() string {
	panic("implement me")
}

func (l *lsetCommand) Type() CommandType {
	return ModifyCommandType
}

func (l *lsetCommand) Keys() []string {
	return []string{l.key}
}

func (l *lsetCommand) ShouldCreate() bool {
	return false
}

func (l *lsetCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	l.accessObject = objects[0]
}

func (l *lsetCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type lremCommand struct {
	key          string
	index        int
	count        int
	item         protocol.RedisString
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (l *lremCommand) Name() string {
	return "lrem"
}

func (l *lremCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 3 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.key = tmpObj.Data()

	countObj, ok := objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	tmpIdx, err := util.ParseInt64(countObj.Data())
	if err != nil {
		return ErrArgumentInvalid
	}

	l.count = int(tmpIdx)

	l.item, ok = objects[2].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	return nil
}

func (l *lremCommand) Execute() {
	if l.accessObject == nil {
		l.result = protocol.NewRedisInteger(0)
		return
	}

	if l.accessObject.Type() != l.TargetContainerType() {
		l.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", l.TargetContainerType(), l.accessObject.Type())
		return
	}

	count := l.accessObject.(container.ListContainer).Remove(l.count, container.NewString(l.item.Data()))
	l.result = protocol.NewRedisInteger(int64(count))
}

func (l *lremCommand) Result() (protocol.RedisObject, error) {
	return l.result, l.err
}

func (l *lremCommand) Cluster() int {
	return l.index
}

func (l *lremCommand) ToLog() string {
	panic("implement me")
}

func (l *lremCommand) Type() CommandType {
	return ModifyCommandType
}

func (l *lremCommand) Keys() []string {
	return []string{l.key}
}

func (l *lremCommand) ShouldCreate() bool {
	return false
}

func (l *lremCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	l.accessObject = objects[0]
}

func (l *lremCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type ltrimCommand struct {
	key          string
	index        int
	start        int
	end          int
	accessObject container.ContainerObject
	result       protocol.RedisString
	err          error
}

func (l *ltrimCommand) Name() string {
	return "ltrim"
}

func (l *ltrimCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 3 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.key = tmpObj.Data()

	startObj, ok := objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	tmpIdx, err := util.ParseInt64(startObj.Data())
	if err != nil {
		return ErrArgumentInvalid
	}

	l.start = int(tmpIdx)

	endObj, ok := objects[2].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	tmpIdx, err = util.ParseInt64(endObj.Data())
	if err != nil {
		return ErrArgumentInvalid
	}

	l.end = int(tmpIdx)

	return nil
}

func (l *ltrimCommand) Execute() {
	if l.accessObject == nil {
		l.result = protocol.NewSimpleRedisString("OK")
		return
	}

	if l.accessObject.Type() != l.TargetContainerType() {
		l.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", l.TargetContainerType(), l.accessObject.Type())
		return
	}

	l.err = l.accessObject.(container.ListContainer).Trim(l.start, l.end)
	l.result = protocol.NewSimpleRedisString("OK")
}

func (l *ltrimCommand) Result() (protocol.RedisObject, error) {
	return l.result, l.err
}

func (l *ltrimCommand) Cluster() int {
	return l.index
}

func (l *ltrimCommand) ToLog() string {
	panic("implement me")
}

func (l *ltrimCommand) Type() CommandType {
	return AccessCommandType
}

func (l *ltrimCommand) Keys() []string {
	return []string{l.key}
}

func (l *ltrimCommand) ShouldCreate() bool {
	return false
}

func (l *ltrimCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	l.accessObject = objects[0]
}

func (l *ltrimCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}

type linsertCommand struct {
	key          string
	index        int
	after        bool
	pivot        protocol.RedisString
	replace      protocol.RedisString
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (l *linsertCommand) Name() string {
	return "linsert"
}

func (l *linsertCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 4 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.key = tmpObj.Data()

	isAfterObj, ok := objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	switch strings.ToLower(isAfterObj.Data()) {
	case "before":
		l.after = false
	case "after":
		l.after = true
	default:
		return ErrArgumentInvalid
	}

	l.pivot, ok = objects[2].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	l.replace, ok = objects[3].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	return nil
}

func (l *linsertCommand) Execute() {
	if l.accessObject == nil {
		l.result = protocol.NewRedisInteger(0)
		return
	}

	if l.accessObject.Type() != l.TargetContainerType() {
		l.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", l.TargetContainerType(), l.accessObject.Type())
		return
	}

	ret, _ := l.accessObject.(container.ListContainer).Insert(container.NewString(l.pivot.Data()), container.NewString(l.replace.Data()), l.after)

	l.result = protocol.NewRedisInteger(int64(ret))
}

func (l *linsertCommand) Result() (protocol.RedisObject, error) {
	return l.result, l.err
}

func (l *linsertCommand) Cluster() int {
	return l.index
}

func (l *linsertCommand) ToLog() string {
	panic("implement me")
}

func (l *linsertCommand) Type() CommandType {
	return ModifyCommandType
}

func (l *linsertCommand) Keys() []string {
	return []string{l.key}
}

func (l *linsertCommand) ShouldCreate() bool {
	return false
}

func (l *linsertCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	l.accessObject = objects[0]
}

func (l *linsertCommand) TargetContainerType() container.ContainerType {
	return container.LinkedListType
}
