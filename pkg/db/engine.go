package db

import "sync"

type Engine interface {
}

type engine struct {
	dbMap sync.Map
}
