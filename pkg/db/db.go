package db

import "github.com/lxdlam/vertex/pkg/container"

type DB interface{

}

type db struct {
	index      int
	containers container.Containers
}

func NewDB(index int) DB {
	return &db {
		index: index,
		containers: container.NewContainers(),
	}
}
