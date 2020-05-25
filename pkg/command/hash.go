package command

import (
	"fmt"

	"github.com/lxdlam/vertex/pkg/container"
	"github.com/lxdlam/vertex/pkg/protocol"
)

func newHashCommand(name string, index int, arguments []protocol.RedisObject) (Command, error) {
	switch name {
	case "hset":
		h := &hsetCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	case "hget":
		h := &hgetCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	case "hmget":
		h := &hmgetCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	case "hexists":
		h := &hexistsCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	case "hlen":
		h := &hlenCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	case "hgetall":
		h := &hgetallCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	case "hkeys":
		h := &hkeysCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	case "hvals":
		h := &hvalsCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	case "hstrlen":
		h := &hstrlenCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	case "hdel":
		h := &hdelCommand{
			index: index,
		}
		err := h.ParseArguments(arguments)
		return h, err
	}
	return nil, ErrCommandNotExist
}

type hsetCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	fields       []string
	values       []string
	result       protocol.RedisInteger
	err          error
}

func (h *hsetCommand) Name() string {
	return "hset"
}

func (h *hsetCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)

	if length%2 != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	for idx := 1; idx < length; idx += 2 {
		keyObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		h.fields = append(h.fields, keyObj.Data())

		valueObj, ok := objects[idx+1].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		h.values = append(h.values, valueObj.Data())
	}

	return nil
}

func (h *hsetCommand) Execute() {
	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	if h.fields == nil || h.values == nil {
		h.err = ErrArgumentInvalid
		return
	}

	var keys, values []*container.StringContainer
	length := len(h.fields)

	for idx := 0; idx < length; idx++ {
		keys = append(keys, container.NewString(h.fields[idx]))
		values = append(values, container.NewString(h.values[idx]))
	}

	ret, _ := h.accessObject.(container.HashContainer).Set(keys, values)

	h.result = protocol.NewRedisInteger(int64(ret))
}

func (h *hsetCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hsetCommand) Cluster() int {
	return h.index
}

func (h *hsetCommand) ToLog() string {
	panic("implement me")
}

func (h *hsetCommand) Type() CommandType {
	return ModifyCommandType
}

func (h *hsetCommand) Keys() []string {
	return []string{h.key}
}

func (h *hsetCommand) ShouldCreate() bool {
	return true
}

func (h *hsetCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hsetCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}

type hgetCommand struct {
	key          string
	index        int
	field        string
	accessObject container.ContainerObject
	result       protocol.RedisString
	err          error
}

func (h *hgetCommand) Name() string {
	return "hget"
}

func (h *hgetCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	tmpObj, ok = objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.field = tmpObj.Data()

	return nil
}

func (h *hgetCommand) Execute() {
	if h.accessObject == nil {
		h.result = protocol.NewNullBulkRedisString()
		return
	}

	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	if h.field == "" {
		h.err = ErrArgumentInvalid
		return
	}

	ret := h.accessObject.(container.HashContainer).Get([]*container.StringContainer{container.NewString(h.field)})

	if len(ret) != 1 || ret[0] == nil {
		h.result = protocol.NewNullBulkRedisString()
	} else {
		h.result = protocol.NewBulkRedisString(ret[0].String())
	}
}

func (h *hgetCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hgetCommand) Cluster() int {
	return h.index
}

func (h *hgetCommand) ToLog() string {
	panic("implement me")
}

func (h *hgetCommand) Type() CommandType {
	return AccessCommandType
}

func (h *hgetCommand) Keys() []string {
	return []string{h.key}
}

func (h *hgetCommand) ShouldCreate() bool {
	return false
}

func (h *hgetCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hgetCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}

type hmgetCommand struct {
	key          string
	index        int
	fields       []string
	accessObject container.ContainerObject
	result       protocol.RedisArray
	err          error
}

func (h *hmgetCommand) Name() string {
	return "hmget"
}

func (h *hmgetCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)

	if length <= 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	for idx := 1; idx < length; idx++ {
		tmpObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		h.fields = append(h.fields, tmpObj.Data())
	}

	return nil
}

func (h *hmgetCommand) Execute() {
	if h.accessObject == nil {
		var objs []protocol.RedisObject

		l := len(h.fields)
		for idx := 0; idx < l; idx++ {
			objs = append(objs, protocol.NewNullBulkRedisString())
		}

		h.result = protocol.NewRedisArray(objs)

		return
	}

	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	var keys []*container.StringContainer

	for _, key := range h.fields {
		keys = append(keys, container.NewString(key))
	}

	var values []protocol.RedisObject

	ret := h.accessObject.(container.HashContainer).Get(keys)
	length := len(ret)

	for idx := 0; idx < length; idx++ {
		if ret[idx] == nil {
			values = append(values, protocol.NewNullBulkRedisString())
		} else {
			values = append(values, protocol.NewBulkRedisString(ret[idx].String()))
		}
	}

	h.result = protocol.NewRedisArray(values)
}

func (h *hmgetCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hmgetCommand) Cluster() int {
	return h.index
}

func (h *hmgetCommand) ToLog() string {
	panic("implement me")
}

func (h *hmgetCommand) Type() CommandType {
	return AccessCommandType
}

func (h *hmgetCommand) Keys() []string {
	return []string{h.key}
}

func (h *hmgetCommand) ShouldCreate() bool {
	return false
}

func (h *hmgetCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hmgetCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}

type hexistsCommand struct {
	key          string
	index        int
	field        string
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (h *hexistsCommand) Name() string {
	return "hexists"
}

func (h *hexistsCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	tmpObj, ok = objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.field = tmpObj.Data()

	return nil
}

func (h *hexistsCommand) Execute() {
	if h.accessObject == nil {
		h.result = protocol.NewRedisInteger(0)
		return
	}

	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	if h.field == "" {
		h.err = ErrArgumentInvalid
		return
	}

	if h.accessObject.(container.HashContainer).Exists(container.NewString(h.field)) {
		h.result = protocol.NewRedisInteger(1)
	} else {
		h.result = protocol.NewRedisInteger(0)
	}
}

func (h *hexistsCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hexistsCommand) Cluster() int {
	return h.index
}

func (h *hexistsCommand) ToLog() string {
	panic("implement me")
}

func (h *hexistsCommand) Type() CommandType {
	return AccessCommandType
}

func (h *hexistsCommand) Keys() []string {
	return []string{h.key}
}

func (h *hexistsCommand) ShouldCreate() bool {
	return false
}

func (h *hexistsCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hexistsCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}

type hlenCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (h *hlenCommand) Name() string {
	return "hlen"
}

func (h *hlenCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	return nil
}

func (h *hlenCommand) Execute() {
	if h.accessObject == nil {
		h.result = protocol.NewRedisInteger(0)
		return
	}

	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	ret := h.accessObject.(container.HashContainer).Len()

	h.result = protocol.NewRedisInteger(int64(ret))
}

func (h *hlenCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hlenCommand) Cluster() int {
	return h.index
}

func (h *hlenCommand) ToLog() string {
	panic("implement me")
}

func (h *hlenCommand) Type() CommandType {
	return AccessCommandType
}

func (h *hlenCommand) Keys() []string {
	return []string{h.key}
}

func (h *hlenCommand) ShouldCreate() bool {
	return false
}

func (h *hlenCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hlenCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}

type hgetallCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisArray
	err          error
}

func (h *hgetallCommand) Name() string {
	return "hgetall"
}

func (h *hgetallCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	return nil
}

func (h *hgetallCommand) Execute() {
	if h.accessObject == nil {
		h.result = protocol.NewNullRedisArray()
		return
	}

	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	var objs []protocol.RedisObject

	fields, values := h.accessObject.(container.HashContainer).Entries()
	l := len(fields)

	for idx := 0; idx < l; idx++ {
		objs = append(objs, protocol.NewBulkRedisString(fields[idx].String()))
		objs = append(objs, protocol.NewBulkRedisString(values[idx].String()))
	}

	h.result = protocol.NewRedisArray(objs)
}

func (h *hgetallCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hgetallCommand) Cluster() int {
	return h.index
}

func (h *hgetallCommand) ToLog() string {
	panic("implement me")
}

func (h *hgetallCommand) Type() CommandType {
	return AccessCommandType
}

func (h *hgetallCommand) Keys() []string {
	return []string{h.key}
}

func (h *hgetallCommand) ShouldCreate() bool {
	return false
}

func (h *hgetallCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hgetallCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}

type hkeysCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisArray
	err          error
}

func (h *hkeysCommand) Name() string {
	return "hkeys"
}

func (h *hkeysCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	return nil
}

func (h *hkeysCommand) Execute() {
	if h.accessObject == nil {
		h.result = protocol.NewNullRedisArray()
		return
	}

	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	var objs []protocol.RedisObject

	fields := h.accessObject.(container.HashContainer).Keys()
	l := len(fields)

	for idx := 0; idx < l; idx++ {
		objs = append(objs, protocol.NewBulkRedisString(fields[idx].String()))
	}

	h.result = protocol.NewRedisArray(objs)
}

func (h *hkeysCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hkeysCommand) Cluster() int {
	return h.index
}

func (h *hkeysCommand) ToLog() string {
	panic("implement me")
}

func (h *hkeysCommand) Type() CommandType {
	return AccessCommandType
}

func (h *hkeysCommand) Keys() []string {
	return []string{h.key}
}

func (h *hkeysCommand) ShouldCreate() bool {
	return false
}

func (h *hkeysCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hkeysCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}

type hvalsCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisArray
	err          error
}

func (h *hvalsCommand) Name() string {
	return "hvals"
}

func (h *hvalsCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	return nil
}

func (h *hvalsCommand) Execute() {
	if h.accessObject == nil {
		h.result = protocol.NewNullRedisArray()
		return
	}

	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	var objs []protocol.RedisObject

	values := h.accessObject.(container.HashContainer).Values()
	l := len(values)

	for idx := 0; idx < l; idx++ {
		objs = append(objs, protocol.NewBulkRedisString(values[idx].String()))
	}

	h.result = protocol.NewRedisArray(objs)
}

func (h *hvalsCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hvalsCommand) Cluster() int {
	return h.index
}

func (h *hvalsCommand) ToLog() string {
	panic("implement me")
}

func (h *hvalsCommand) Type() CommandType {
	return AccessCommandType
}

func (h *hvalsCommand) Keys() []string {
	return []string{h.key}
}

func (h *hvalsCommand) ShouldCreate() bool {
	return false
}

func (h *hvalsCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hvalsCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}

type hstrlenCommand struct {
	key          string
	index        int
	field        string
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (h *hstrlenCommand) Name() string {
	return "hstrlen"
}

func (h *hstrlenCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	tmpObj, ok = objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.field = tmpObj.Data()

	return nil
}

func (h *hstrlenCommand) Execute() {
	if h.accessObject == nil {
		h.result = protocol.NewRedisInteger(0)
		return
	}

	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	if h.field == "" {
		h.err = ErrArgumentInvalid
		return
	}

	ret := h.accessObject.(container.HashContainer).Get([]*container.StringContainer{container.NewString(h.field)})

	if len(ret) != 1 || ret[0] == nil {
		h.result = protocol.NewRedisInteger(0)
	} else {
		h.result = protocol.NewRedisInteger(int64(ret[0].Len()))
	}
}

func (h *hstrlenCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hstrlenCommand) Cluster() int {
	return h.index
}

func (h *hstrlenCommand) ToLog() string {
	panic("implement me")
}

func (h *hstrlenCommand) Type() CommandType {
	return AccessCommandType
}

func (h *hstrlenCommand) Keys() []string {
	return []string{h.key}
}

func (h *hstrlenCommand) ShouldCreate() bool {
	return false
}

func (h *hstrlenCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hstrlenCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}

type hdelCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	fields       []string
	result       protocol.RedisInteger
	err          error
}

func (h *hdelCommand) Name() string {
	return "hdel"
}

func (h *hdelCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)

	if length <= 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	h.key = tmpObj.Data()

	for idx := 1; idx < length; idx++ {
		keyObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		h.fields = append(h.fields, keyObj.Data())
	}

	return nil
}

func (h *hdelCommand) Execute() {
	if h.accessObject == nil {
		h.result = protocol.NewRedisInteger(0)
		return
	}

	if h.accessObject.Type() != h.TargetContainerType() {
		h.err = fmt.Errorf("target container type mismatch. expected=%d, got=%d", h.TargetContainerType(), h.accessObject.Type())
		return
	}

	if len(h.fields) == 0 {
		h.err = ErrArgumentInvalid
		return
	}

	var fields []*container.StringContainer
	length := len(h.fields)

	for idx := 0; idx < length; idx++ {
		fields = append(fields, container.NewString(h.fields[idx]))
	}

	ret := h.accessObject.(container.HashContainer).Del(fields)

	h.result = protocol.NewRedisInteger(int64(ret))
}

func (h *hdelCommand) Result() (protocol.RedisObject, error) {
	return h.result, h.err
}

func (h *hdelCommand) Cluster() int {
	return h.index
}

func (h *hdelCommand) ToLog() string {
	panic("implement me")
}

func (h *hdelCommand) Type() CommandType {
	return ModifyCommandType
}

func (h *hdelCommand) Keys() []string {
	return []string{h.key}
}

func (h *hdelCommand) ShouldCreate() bool {
	return false
}

func (h *hdelCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	h.accessObject = objects[0]
}

func (h *hdelCommand) TargetContainerType() container.ContainerType {
	return container.HashType
}
