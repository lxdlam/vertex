package util

import "github.com/lxdlam/vertex/pkg/types"

// Handler is type alias for better readability
type Handler = func(types.DataMap) (types.DataMap, error)

// Pipeline is a simple converted to package a series of works
// that we may need to check whether the latter step is succeeded
// to process.
type Pipeline interface {
	AddHandler(Handler) Pipeline

	Run(types.DataMap) (types.DataMap, error)
}

type pipeline struct {
	handlers []Handler
}

// NewPipeline returns a new pipeline instance
func NewPipeline() Pipeline {
	return &pipeline{handlers: nil}
}

func (p *pipeline) AddHandler(handler Handler) Pipeline {
	p.handlers = append(p.handlers, handler)

	// for chain call
	return p
}

func (p *pipeline) Run(ctx types.DataMap) (types.DataMap, error) {
	var err error = nil
	for _, handler := range p.handlers {
		ctx, err = handler(ctx)
		if err != nil {
			break
		}
	}

	return ctx, err
}
