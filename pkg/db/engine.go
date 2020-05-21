package db

import "sync"

type Engine interface {
}

// TODO(2020.5.22): currently only one db instance is supported, regardless the value required by client
type engine struct {
	dbMap    sync.Map
	shutChan chan struct{}
}

// NewEngine will return a new engine that listens to the requests and post responses
func NewEngine() Engine {
	e := &engine{}

	// TODO(2020.5.22): create a new instance
	e.dbMap.Store(1, NewDB(1))

	// TODO: insert clean codes
	go func() {
	Outer:
		for {
			select {
			case <-e.shutChan:
				break Outer
			}
		}
	}()

	return e
}
