package db

import (
	"github.com/lxdlam/vertex/pkg/command"
	"github.com/lxdlam/vertex/pkg/container"
)

type DB interface {
	ExecuteCommand(command.Command)

	Index() int
}

type db struct {
	index      int
	containers container.Containers
	shutChan   chan struct{}
}

func NewDB(index int) DB {
	return &db{
		index:      index,
		containers: container.NewContainers(),
	}
}

func (d *db) resolveName(key string, t container.ContainerType, create bool) container.ContainerObject {
	switch t {
	case container.GlobalType:
		// create do nothing here
		return d.containers.Global()
	case container.LinkedListType:
		if create {
			return d.containers.GetOrCreateList(key)
		} else {
			return d.containers.GetList(key)
		}
	case container.HashType:
		if create {
			return d.containers.GetOrCreateHash(key)
		} else {
			return d.containers.GetHash(key)
		}
	case container.SetType:
		if create {
			return d.containers.GetOrCreateSet(key)
		} else {
			return d.containers.GetSet(key)
		}
	}

	return nil
}

func (d *db) ExecuteCommand(c command.Command) {
	var accessObjects []container.ContainerObject
	targetType := c.TargetContainerType()
	create := c.ShouldCreate()

	for _, key := range c.Keys() {
		accessObjects = append(accessObjects, d.resolveName(key, targetType, create))
	}

	c.SetAccessObjects(accessObjects)

	c.Execute()
}

func (d *db) Index() int {
	return d.index
}
