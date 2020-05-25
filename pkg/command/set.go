package command

import (
	"fmt"

	"github.com/lxdlam/vertex/pkg/container"
	"github.com/lxdlam/vertex/pkg/protocol"
	"github.com/lxdlam/vertex/pkg/util"
)

func newSetCommand(name string, index int, arguments []protocol.RedisObject) (Command, error) {
	switch name {
	case "scard":
		s := &scardCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "smembers":
		s := &smembersCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "sadd":
		s := &saddCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "sismember":
		s := &sismemberCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "srem":
		s := &sremCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "spop":
		s := &spopCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "srandmember":
		s := &srandmemberCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "sdiff":
		s := &sdiffCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "sinter":
		s := &sinterCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	case "sunion":
		s := &sunionCommand{
			index: index,
		}
		err := s.ParseArguments(arguments)
		return s, err
	}
	return nil, ErrCommandNotExist
}

type scardCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (s *scardCommand) Name() string {
	return "scard"
}

func (s *scardCommand) ParseArguments(objects []protocol.RedisObject) error {
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

func (s *scardCommand) Execute() {
	if s.accessObject == nil {
		s.result = protocol.NewRedisInteger(0)
		return
	}

	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	ret := s.accessObject.(container.SetContainer).Len()

	s.result = protocol.NewRedisInteger(int64(ret))
}

func (s *scardCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *scardCommand) Cluster() int {
	return s.index
}

func (s *scardCommand) ToLog() string {
	panic("implement me")
}

func (s *scardCommand) Type() CommandType {
	return AccessCommandType
}

func (s *scardCommand) Keys() []string {
	return []string{s.key}
}

func (s *scardCommand) ShouldCreate() bool {
	return false
}

func (s *scardCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *scardCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}

type saddCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	values       []string
	result       protocol.RedisInteger
	err          error
}

func (s *saddCommand) Name() string {
	return "sadd"
}

func (s *saddCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)

	if length <= 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	s.key = tmpObj.Data()

	for idx := 1; idx < length; idx++ {
		tmpObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		s.values = append(s.values, tmpObj.Data())
	}

	return nil
}

func (s *saddCommand) Execute() {
	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	if len(s.values) == 0 {
		s.err = ErrArgumentInvalid
		return
	}

	var values []*container.StringContainer
	length := len(s.values)

	for idx := 0; idx < length; idx++ {
		values = append(values, container.NewString(s.values[idx]))
	}

	ret := s.accessObject.(container.SetContainer).Add(values)

	s.result = protocol.NewRedisInteger(int64(ret))
}

func (s *saddCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *saddCommand) Cluster() int {
	return s.index
}

func (s *saddCommand) ToLog() string {
	panic("implement me")
}

func (s *saddCommand) Type() CommandType {
	return ModifyCommandType
}

func (s *saddCommand) Keys() []string {
	return []string{s.key}
}

func (s *saddCommand) ShouldCreate() bool {
	return true
}

func (s *saddCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *saddCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}

type smembersCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	result       protocol.RedisArray
	err          error
}

func (s *smembersCommand) Name() string {
	return "smembers"
}

func (s *smembersCommand) ParseArguments(objects []protocol.RedisObject) error {
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

func (s *smembersCommand) Execute() {
	if s.accessObject == nil {
		s.result = protocol.NewNullRedisArray()
		return
	}

	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	var objs []protocol.RedisObject

	fields := s.accessObject.(container.SetContainer).Members()
	l := len(fields)

	for idx := 0; idx < l; idx++ {
		objs = append(objs, protocol.NewBulkRedisString(fields[idx].String()))
	}

	s.result = protocol.NewRedisArray(objs)
}

func (s *smembersCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *smembersCommand) Cluster() int {
	return s.index
}

func (s *smembersCommand) ToLog() string {
	panic("implement me")
}

func (s *smembersCommand) Type() CommandType {
	return AccessCommandType
}

func (s *smembersCommand) Keys() []string {
	return []string{s.key}
}

func (s *smembersCommand) ShouldCreate() bool {
	return false
}

func (s *smembersCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *smembersCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}

type sismemberCommand struct {
	key          string
	index        int
	value        string
	accessObject container.ContainerObject
	result       protocol.RedisInteger
	err          error
}

func (s *sismemberCommand) Name() string {
	return "sismember"
}

func (s *sismemberCommand) ParseArguments(objects []protocol.RedisObject) error {
	if len(objects) != 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	s.key = tmpObj.Data()

	tmpObj, ok = objects[1].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	s.value = tmpObj.Data()

	return nil
}

func (s *sismemberCommand) Execute() {
	if s.accessObject == nil {
		s.result = protocol.NewRedisInteger(0)
		return
	}

	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	if s.value == "" {
		s.err = ErrArgumentInvalid
		return
	}

	if s.accessObject.(container.SetContainer).IsMember(container.NewString(s.value)) {
		s.result = protocol.NewRedisInteger(1)
	} else {
		s.result = protocol.NewRedisInteger(0)
	}
}

func (s *sismemberCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *sismemberCommand) Cluster() int {
	return s.index
}

func (s *sismemberCommand) ToLog() string {
	panic("implement me")
}

func (s *sismemberCommand) Type() CommandType {
	return AccessCommandType
}

func (s *sismemberCommand) Keys() []string {
	return []string{s.key}
}

func (s *sismemberCommand) ShouldCreate() bool {
	return false
}

func (s *sismemberCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *sismemberCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}

type sremCommand struct {
	key          string
	index        int
	accessObject container.ContainerObject
	values       []string
	result       protocol.RedisInteger
	err          error
}

func (s *sremCommand) Name() string {
	return "srem"
}

func (s *sremCommand) ParseArguments(objects []protocol.RedisObject) error {
	length := len(objects)

	if length <= 1 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	s.key = tmpObj.Data()

	for idx := 1; idx < length; idx++ {
		tmpObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		s.values = append(s.values, tmpObj.Data())
	}

	return nil
}

func (s *sremCommand) Execute() {
	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	if len(s.values) == 0 {
		s.err = ErrArgumentInvalid
		return
	}

	var values []*container.StringContainer
	length := len(s.values)

	for idx := 0; idx < length; idx++ {
		values = append(values, container.NewString(s.values[idx]))
	}

	ret := s.accessObject.(container.SetContainer).Delete(values)

	s.result = protocol.NewRedisInteger(int64(ret))
}

func (s *sremCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *sremCommand) Cluster() int {
	return s.index
}

func (s *sremCommand) ToLog() string {
	panic("implement me")
}

func (s *sremCommand) Type() CommandType {
	return ModifyCommandType
}

func (s *sremCommand) Keys() []string {
	return []string{s.key}
}

func (s *sremCommand) ShouldCreate() bool {
	return true
}

func (s *sremCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *sremCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}

type spopCommand struct {
	key          string
	index        int
	count        int
	countSet     bool
	accessObject container.ContainerObject
	result       protocol.RedisObject
	err          error
}

func (s *spopCommand) Name() string {
	return "spop"
}

func (s *spopCommand) ParseArguments(objects []protocol.RedisObject) error {
	l := len(objects)

	if l <= 0 || l > 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	s.key = tmpObj.Data()

	if l == 2 {
		countObj, ok := objects[1].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		tmpIdx, err := util.ParseInt64(countObj.Data())
		if err != nil {
			return ErrArgumentInvalid
		}

		s.count = int(tmpIdx)
		s.countSet = true
	} else {
		s.count = 1
		s.countSet = false
	}

	return nil
}

func (s *spopCommand) Execute() {
	if s.accessObject == nil {
		if !s.countSet {
			s.result = protocol.NewNullBulkRedisString()
		} else {
			s.result = protocol.NewNullRedisArray()
		}
		return
	}

	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	ret := s.accessObject.(container.SetContainer).Pop(s.count)

	if !s.countSet {
		if ret[0] == nil {
			s.result = protocol.NewNullBulkRedisString()
		} else {
			s.result = protocol.NewBulkRedisString(ret[0].String())
		}
	} else {
		if len(ret) == 0 {
			s.result = protocol.NewNullRedisArray()
		} else {
			var objs []protocol.RedisObject

			for _, obj := range ret {
				objs = append(objs, protocol.NewBulkRedisString(obj.String()))
			}

			s.result = protocol.NewRedisArray(objs)
		}
	}
}

func (s *spopCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *spopCommand) Cluster() int {
	return s.index
}

func (s *spopCommand) ToLog() string {
	panic("implement me")
}

func (s *spopCommand) Type() CommandType {
	return ModifyCommandType
}

func (s *spopCommand) Keys() []string {
	return []string{s.key}
}

func (s *spopCommand) ShouldCreate() bool {
	return false
}

func (s *spopCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *spopCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}

type srandmemberCommand struct {
	key          string
	index        int
	count        int
	countSet     bool
	accessObject container.ContainerObject
	result       protocol.RedisObject
	err          error
}

func (s *srandmemberCommand) Name() string {
	return "srandmember"
}

func (s *srandmemberCommand) ParseArguments(objects []protocol.RedisObject) error {
	l := len(objects)

	if l <= 0 || l > 2 {
		return ErrArgumentInvalid
	}

	tmpObj, ok := objects[0].(protocol.RedisString)
	if !ok {
		return ErrArgumentInvalid
	}

	s.key = tmpObj.Data()

	if l == 2 {
		countObj, ok := objects[1].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		tmpIdx, err := util.ParseInt64(countObj.Data())
		if err != nil {
			return ErrArgumentInvalid
		}

		s.count = int(tmpIdx)
		s.countSet = true
	} else {
		s.count = 1
		s.countSet = false
	}

	return nil
}

func (s *srandmemberCommand) Execute() {
	if s.accessObject == nil {
		if !s.countSet {
			s.result = protocol.NewNullBulkRedisString()
		} else {
			s.result = protocol.NewNullRedisArray()
		}
		return
	}

	if s.accessObject.Type() != s.TargetContainerType() {
		s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), s.accessObject.Type())
		return
	}

	ret := s.accessObject.(container.SetContainer).RandomMember(s.count)

	if !s.countSet {
		if ret[0] == nil {
			s.result = protocol.NewNullBulkRedisString()
		} else {
			s.result = protocol.NewBulkRedisString(ret[0].String())
		}
	} else {
		if len(ret) == 0 {
			s.result = protocol.NewNullRedisArray()
		} else {
			var objs []protocol.RedisObject

			for _, obj := range ret {
				objs = append(objs, protocol.NewBulkRedisString(obj.String()))
			}

			s.result = protocol.NewRedisArray(objs)
		}
	}
}

func (s *srandmemberCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *srandmemberCommand) Cluster() int {
	return s.index
}

func (s *srandmemberCommand) ToLog() string {
	panic("implement me")
}

func (s *srandmemberCommand) Type() CommandType {
	return AccessCommandType
}

func (s *srandmemberCommand) Keys() []string {
	return []string{s.key}
}

func (s *srandmemberCommand) ShouldCreate() bool {
	return false
}

func (s *srandmemberCommand) SetAccessObjects(objects []container.ContainerObject) {
	if len(objects) == 0 {
		return
	}
	s.accessObject = objects[0]
}

func (s *srandmemberCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}

type sdiffCommand struct {
	keys          []string
	index         int
	accessObjects []container.ContainerObject
	result        protocol.RedisObject
	err           error
}

func (s *sdiffCommand) Name() string {
	return "sdiff"
}

func (s *sdiffCommand) ParseArguments(objects []protocol.RedisObject) error {
	l := len(objects)

	if l <= 0 {
		return ErrArgumentInvalid
	}

	s.keys = make([]string, l)

	for idx := 0; idx < l; idx++ {
		tmpObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		s.keys[idx] = tmpObj.Data()
	}

	return nil
}

func (s *sdiffCommand) Execute() {
	if s.accessObjects == nil {
		s.result = protocol.NewNullRedisArray()
		return
	}

	for _, accessObject := range s.accessObjects {
		if accessObject.Type() != s.TargetContainerType() {
			s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), accessObject.Type())
			return
		}
	}

	var others []container.SetContainer

	base := s.accessObjects[0].(container.SetContainer)

	for _, item := range s.accessObjects[1:] {
		others = append(others, item.(container.SetContainer))
	}

	ret := base.Diff(others)

	if ret.Len() == 0 {
		s.result = protocol.NewNullRedisArray()
	} else {
		var objs []protocol.RedisObject

		fields := ret.Members()
		l := len(fields)

		for idx := 0; idx < l; idx++ {
			objs = append(objs, protocol.NewBulkRedisString(fields[idx].String()))
		}

		s.result = protocol.NewRedisArray(objs)
	}
}

func (s *sdiffCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *sdiffCommand) Cluster() int {
	return s.index
}

func (s *sdiffCommand) ToLog() string {
	panic("implement me")
}

func (s *sdiffCommand) Type() CommandType {
	return AccessCommandType
}

func (s *sdiffCommand) Keys() []string {
	return s.keys
}

func (s *sdiffCommand) ShouldCreate() bool {
	return false
}

func (s *sdiffCommand) SetAccessObjects(objects []container.ContainerObject) {
	l := len(objects)
	s.accessObjects = make([]container.ContainerObject, l)

	for idx := 0; idx < l; idx++ {
		if objects[idx] == nil {
			s.accessObjects[idx] = container.NewSetContainer("anonymous")
		} else {
			s.accessObjects[idx] = objects[idx]
		}
	}
}

func (s *sdiffCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}

type sinterCommand struct {
	keys          []string
	index         int
	accessObjects []container.ContainerObject
	result        protocol.RedisObject
	err           error
}

func (s *sinterCommand) Name() string {
	return "sinter"
}

func (s *sinterCommand) ParseArguments(objects []protocol.RedisObject) error {
	l := len(objects)

	if l <= 0 {
		return ErrArgumentInvalid
	}

	s.keys = make([]string, l)

	for idx := 0; idx < l; idx++ {
		tmpObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		s.keys[idx] = tmpObj.Data()
	}

	return nil
}

func (s *sinterCommand) Execute() {
	if s.accessObjects == nil {
		s.result = protocol.NewNullRedisArray()
		return
	}

	for _, accessObject := range s.accessObjects {
		if accessObject.Type() != s.TargetContainerType() {
			s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), accessObject.Type())
			return
		}
	}

	var others []container.SetContainer

	base := s.accessObjects[0].(container.SetContainer)

	for _, item := range s.accessObjects[1:] {
		others = append(others, item.(container.SetContainer))
	}

	ret := base.Intersect(others)

	if ret.Len() == 0 {
		s.result = protocol.NewNullRedisArray()
	} else {
		var objs []protocol.RedisObject

		fields := ret.Members()
		l := len(fields)

		for idx := 0; idx < l; idx++ {
			objs = append(objs, protocol.NewBulkRedisString(fields[idx].String()))
		}

		s.result = protocol.NewRedisArray(objs)
	}
}

func (s *sinterCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *sinterCommand) Cluster() int {
	return s.index
}

func (s *sinterCommand) ToLog() string {
	panic("implement me")
}

func (s *sinterCommand) Type() CommandType {
	return AccessCommandType
}

func (s *sinterCommand) Keys() []string {
	return s.keys
}

func (s *sinterCommand) ShouldCreate() bool {
	return false
}

func (s *sinterCommand) SetAccessObjects(objects []container.ContainerObject) {
	l := len(objects)
	s.accessObjects = make([]container.ContainerObject, l)

	for idx := 0; idx < l; idx++ {
		if objects[idx] == nil {
			s.accessObjects[idx] = container.NewSetContainer("anonymous")
		} else {
			s.accessObjects[idx] = objects[idx]
		}
	}
}

func (s *sinterCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}

type sunionCommand struct {
	keys          []string
	index         int
	accessObjects []container.ContainerObject
	result        protocol.RedisObject
	err           error
}

func (s *sunionCommand) Name() string {
	return "sunion"
}

func (s *sunionCommand) ParseArguments(objects []protocol.RedisObject) error {
	l := len(objects)

	if l <= 0 {
		return ErrArgumentInvalid
	}

	s.keys = make([]string, l)

	for idx := 0; idx < l; idx++ {
		tmpObj, ok := objects[idx].(protocol.RedisString)
		if !ok {
			return ErrArgumentInvalid
		}

		s.keys[idx] = tmpObj.Data()
	}

	return nil
}

func (s *sunionCommand) Execute() {
	if s.accessObjects == nil {
		s.result = protocol.NewNullRedisArray()
		return
	}

	for _, accessObject := range s.accessObjects {
		if accessObject.Type() != s.TargetContainerType() {
			s.err = fmt.Errorf("target container type mismatcs. expected=%d, got=%d", s.TargetContainerType(), accessObject.Type())
			return
		}
	}

	var others []container.SetContainer

	base := s.accessObjects[0].(container.SetContainer)

	for _, item := range s.accessObjects[1:] {
		others = append(others, item.(container.SetContainer))
	}

	ret := base.Union(others)

	if ret.Len() == 0 {
		s.result = protocol.NewNullRedisArray()
	} else {
		var objs []protocol.RedisObject

		fields := ret.Members()
		l := len(fields)

		for idx := 0; idx < l; idx++ {
			objs = append(objs, protocol.NewBulkRedisString(fields[idx].String()))
		}

		s.result = protocol.NewRedisArray(objs)
	}
}

func (s *sunionCommand) Result() (protocol.RedisObject, error) {
	return s.result, s.err
}

func (s *sunionCommand) Cluster() int {
	return s.index
}

func (s *sunionCommand) ToLog() string {
	panic("implement me")
}

func (s *sunionCommand) Type() CommandType {
	return AccessCommandType
}

func (s *sunionCommand) Keys() []string {
	return s.keys
}

func (s *sunionCommand) ShouldCreate() bool {
	return false
}

func (s *sunionCommand) SetAccessObjects(objects []container.ContainerObject) {
	l := len(objects)
	s.accessObjects = make([]container.ContainerObject, l)

	for idx := 0; idx < l; idx++ {
		if objects[idx] == nil {
			s.accessObjects[idx] = container.NewSetContainer("anonymous")
		} else {
			s.accessObjects[idx] = objects[idx]
		}
	}
}

func (s *sunionCommand) TargetContainerType() container.ContainerType {
	return container.SetType
}
