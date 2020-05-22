package util_test

import (
	"errors"
	"testing"

	"github.com/lxdlam/vertex/pkg/types"
	. "github.com/lxdlam/vertex/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestPipelineSuccess(t *testing.T) {
	p := NewPipeline()
	ctx := types.NewSimpleDataMap()
	ctx.Set("number", 0)

	p.AddHandler(simpleAdd).AddHandler(simpleAdd).AddHandler(simpleAdd).AddHandler(simpleAdd).AddHandler(simpleAdd)

	newCtx, err := p.Run(ctx)
	assert.Nil(t, err)

	r, ok := newCtx.Get("number")
	ret := r.(int)

	assert.True(t, ok)
	assert.Equal(t, ret, 5)
}

func TestPipelineFail(t *testing.T) {
	p := NewPipeline()
	ctx := types.NewSimpleDataMap()
	ctx.Set("number", 0)

	p.AddHandler(simpleAdd).AddHandler(simpleAdd).AddHandler(simpleAdd).AddHandler(func(ctx types.DataMap) (types.DataMap, error) {
		return ctx, errors.New("stop here")
	}).AddHandler(simpleAdd)

	newCtx, err := p.Run(ctx)
	assert.Error(t, err, "stop here")

	r, ok := newCtx.Get("number")
	ret := r.(int)

	assert.True(t, ok)
	assert.Equal(t, ret, 3)
}

func simpleAdd(ctx types.DataMap) (types.DataMap, error) {
	num, _ := ctx.Get("number")
	realNum := num.(int)
	realNum++

	ctx.Set("number", realNum)

	return ctx, nil
}
