package db

import (
	"errors"
	"fmt"
	"github.com/lxdlam/vertex/pkg/container"
	"reflect"
	"sync"

	"github.com/lxdlam/vertex/pkg/command"
	"github.com/lxdlam/vertex/pkg/common"
	"github.com/lxdlam/vertex/pkg/concurrency"
	"github.com/lxdlam/vertex/pkg/protocol"
	"github.com/lxdlam/vertex/pkg/types"
)

type Engine interface {
	Start()
	Stop()
}

// TODO(2020.5.22): currently only one db instance is supported, regardless the value required by client
type engine struct {
	dbMap           sync.Map
	requestReceiver concurrency.Receiver
	eventBus        concurrency.EventBus
	shutChan        chan struct{}
}

// NewEngine will return a new engine that listens to the requests and post responses
func NewEngine() Engine {
	e := &engine{
		shutChan: make(chan struct{}),
		eventBus: concurrency.GetEventBus(),
	}

	var err error
	e.requestReceiver, err = e.eventBus.Subscribe("request", "engine")
	if err != nil {
		_ = common.Errorf("subscribe to request failed. err={%w}", err)
		return nil
	}

	return e
}

func (e *engine) getDB(index int) DB {
	db, ok := e.dbMap.Load(index)
	if !ok {
		return nil
	}

	return db.(DB)
}

func (e *engine) getOrCreateDB(index int) DB {
	db, _ := e.dbMap.LoadOrStore(index, NewDB(index))

	return db.(DB)
}

func parseRequest(data types.DataMap) (id string, objects []protocol.RedisObject, ok bool) {
	i, ok := data.Get("id")
	if !ok {
		return
	}

	id, ok = i.(string)
	if !ok {
		return
	}

	response, ok := data.Get("request")
	if !ok {
		return
	}

	obj, ok := response.(protocol.RedisArray)
	if !ok {
		return
	}

	objects = obj.Data()
	return
}

func (e *engine) handleRequest(objects []protocol.RedisObject) (protocol.RedisObject, error) {
	if len(objects) == 0 {
		return nil, fmt.Errorf("empty request objects")
	}

	n, ok := objects[0].(protocol.RedisString)
	if !ok {
		return nil, fmt.Errorf("invalid command name, raw=%+v", objects[0])
	}

	name := n.Data()

	c, err := command.NewCommand(name, 1, objects[1:])

	if err != nil || c == nil {
		return nil, fmt.Errorf("new commond error, name=%s, index=1, error={%w}", name, err)
	}

	db := e.getOrCreateDB(1)
	db.ExecuteCommand(c)

	ret, err := c.Result()
	if err != nil {
		return nil, fmt.Errorf("execute error. name=%s, command=%+v, err={%w}", name, c, err)
	}

	return ret, nil
}

func (e *engine) Start() {
Outer:
	for {
		select {
		case <-e.shutChan:
			break Outer
		default:
			event, err := e.requestReceiver.Receive()

			if err != nil {
				if errors.Is(err, concurrency.ErrChannelClosed) {
					break Outer
				} else {
					common.Warnf("unexpected error. event_id=%s, err=%+v", event.ID(), err)
					break
				}
			}

			if err := event.Error(); err != nil {
				common.Warnf("received and error when receive a new event. event_id=%s, err=%+v", event.ID(), err)
				break
			}

			data, ok := event.Data().(types.DataMap)
			if !ok {
				_ = common.Errorf("data from event is not an DataMap. event_id=%s, raw_type=%s", event.ID(), reflect.TypeOf(data).String())
				break
			}

			id, objects, ok := parseRequest(data)

			if !ok {
				_ = common.Errorf("parse data from DataMap failed. event_id=%s, raw_data=%+v", event.ID(), data)
				break
			}

			responseMap := types.NewSimpleDataMap()
			responseMap.Set("id", id)

			ret, err := e.handleRequest(objects)

			if err != nil {
				responseMap.Set("response", handleError(err))
			} else {
				responseMap.Set("response", ret)
			}

			e.eventBus.Publish("response", responseMap, err)
		}
	}
}

func (e *engine) Stop() {
	close(e.shutChan)
}

func handleError(err error) protocol.RedisError {
	common.Debugf("engine handle request error. err={%+v}", err)
	if errors.Is(err, command.ErrCommandNotExist) {
		return protocol.NewRedisError("ERR no such command")
	} else if errors.Is(err, command.ErrArgumentInvalid) {
		return protocol.NewRedisError("ERR invalid argument")
	} else if errors.Is(err, container.ErrNotAInt) {
		return protocol.NewRedisError("ERR value is not an integer or out of range")
	} else if errors.Is(err, command.ErrNoSuchKey){
		return protocol.NewRedisError("ERR no such key")
	} else if errors.Is(err, container.ErrOutOfRange) {
		return protocol.NewRedisError("ERR index out of range")
	}

	// TODO: do not send raw error
	return protocol.NewRedisError(fmt.Sprintf("ERR vertex server internal error, err=%+v", err))

}
