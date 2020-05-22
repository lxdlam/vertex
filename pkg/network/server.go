package network

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lxdlam/vertex/pkg/types"

	"github.com/lxdlam/vertex/pkg/common"
	"github.com/lxdlam/vertex/pkg/concurrency"
	"github.com/lxdlam/vertex/pkg/protocol"
)

var (
	// FatalResponse will be returned if we met any fatal error when we answering the client
	FatalResponse = protocol.NewRedisError("fatal error by vertex. check server log.")
)

// Server is the main service of vertex.
type Server interface {
	// Prepare the resource, false if failed.
	// If it returns false, the server should not be start.
	Init(c common.Config) bool

	// Serve call will start the run session with blocking the current goroutine.
	// If it returns, it means that the server is requested to be stopped.
	Serve()

	// Stop will stop the server. It provides a way to manually stop the server.
	Stop()
}

// server load will be high, so we use sync.Map, the memory overhead should be profiled
type server struct {
	tcpListener      *net.TCPListener
	addr             *net.TCPAddr
	eventBus         concurrency.EventBus
	shutChan         chan os.Signal
	shutdown         int32
	responseReceiver concurrency.Receiver
	cleanUpHandles   []func()
	clients          sync.Map
}

// NewServer will returns a new server instance
func NewServer() Server {
	return &server{
		tcpListener:      nil,
		eventBus:         concurrency.GetEventBus(),
		shutChan:         nil,
		shutdown:         0,
		responseReceiver: nil,
		cleanUpHandles:   nil,
		clients:          sync.Map{},
	}
}

func (s *server) Init(c common.Config) bool {
	var err error

	common.Debug("initialize server")
	defer func() {
		if err == nil {
			common.Debug("initialize server success")
		} else {
			common.Fatalf("initialize server error. err=%s", err.Error())
		}
	}()

	s.addr = &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: c.Port,
	}

	s.tcpListener, err = net.ListenTCP("tcp", s.addr)
	if err != nil {
		_ = common.Errorf("init tcp listener failed. addr=%+v, err={%s}", s.addr, err.Error())

		return false
	}

	s.eventBus.NewTopic("request")
	s.eventBus.NewTopic("response")

	// Set buffer size to 100
	s.responseReceiver, err = s.eventBus.SubscribeWithOptions("response", "server", 100, 10*time.Millisecond)
	if err != nil {
		common.Fatalf("Init server in subscribe response channel failed. err=%s", err)
	}

	s.shutChan = make(chan os.Signal)
	signal.Notify(s.shutChan, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGKILL)

	return true
}

func (s *server) Serve() {
	_, _ = fmt.Fprintf(os.Stderr, "Server is listening on %s\n", s.addr.String())
	common.Infof("Server is listening on %s", s.addr.String())

	var wg sync.WaitGroup
	wg.Add(2)

	// Currently we just spawn one goroutine
	// If it's needed, add more workers
	go func() {
	Outer:
		for {
			select {
			case <-s.shutChan:
				break Outer
			default:
				conn, err := s.tcpListener.Accept()
				if err != nil {
					common.Warnf("tcp listen error. err=%s", err)
				}

				// Start handle new conn
				go s.newConn(conn)
			}
		}

		wg.Done()
	}()

	go func() {
		s.responseWorker()
		wg.Done()
	}()

	// wait until shutdown
	wg.Wait()

	s.clean()
}

func (s *server) Stop() {
	if atomic.CompareAndSwapInt32(&s.shutdown, 0, 1) {
		s.shutChan <- syscall.SIGTERM
	}
}

func (s *server) newConn(conn net.Conn) {
	c := NewConn(conn)
	s.clients.Store(c.ID(), c)

	go func() {
		for {
			request, err := c.Read()
			if err != nil {
				common.Warnf("server: new conn with error. addr=%s, err=%s", c.Addr(), err.Error())
				return
			}

			dataMap := types.NewSimpleDataMap()
			dataMap.Set("id", c.ID())
			dataMap.Set("request", request)

			concurrency.GetEventBus().Publish("request", dataMap, nil)
		}
	}()
}

func parseResponse(data types.DataMap) (id string, obj protocol.RedisObject, ok bool) {
	i, ok := data.Get("id")
	if !ok {
		return
	}

	id, ok = i.(string)
	if !ok {
		return
	}

	response, ok := data.Get("response")
	if !ok {
		return
	}

	obj, ok = response.(protocol.RedisObject)
	return
}

func (s *server) responseWorker() {
Outer:
	for {
		select {
		case <-s.shutChan:
			break Outer
		default:
			event, err := s.responseReceiver.Receive()

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

			id, obj, ok := parseResponse(data)

			if !ok {
				_ = common.Errorf("parse data from DataMap failed. event_id=%s, raw_data=%+v", event.ID(), data)
				break
			}

			conn, ok := s.clients.Load(id)

			if ok {
				if err := conn.(Conn).Write(obj.String()); err != nil {
					_ = common.Errorf("send response to client failed. client.addr=%s, object=%s", conn.(Conn).Addr(), obj.String())
					break
				}
			} else {
				_ = common.Errorf("conn not exist. id=%s", id)
			}
		}
	}
}

func (s *server) clean() {
	s.eventBus.RemoveTopic("request")
	s.eventBus.RemoveTopic("response")

	close(s.shutChan)

	s.clients.Range(func(key, value interface{}) bool {
		c, ok := value.(conn)
		if ok {
			_ = c.Close()
		}
		return true
	})
}
