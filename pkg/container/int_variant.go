package container

import (
	"fmt"
	"strconv"

	"github.com/lxdlam/vertex/pkg/protocol"
)

type intVariant interface {
	Get() int64
	Increase(int64)
	Decrease(int64)

	AsString() string
	AsIntObject() protocol.RedisInteger
}

type intVariantImpl struct {
	val int64
}

func newIntVariant(data string) intVariant {
	if val, err := strconv.ParseInt(data, 10, 64); err == nil {
		return &intVariantImpl{val: val}
	}

	return nil
}

func (i *intVariantImpl) Get() int64 {
	return i.val
}

func (i *intVariantImpl) Increase(increment int64) {
	i.val += increment
}

func (i *intVariantImpl) Decrease(decrement int64) {
	i.val -= decrement
}

func (i *intVariantImpl) AsString() string {
	return fmt.Sprintf("%d", i.val)
}

func (i *intVariantImpl) AsIntObject() protocol.RedisInteger {
	return protocol.NewRedisInteger(i.val)
}
