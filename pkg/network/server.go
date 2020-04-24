package network

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lxdlam/vertex/pkg/common"
	"github.com/lxdlam/vertex/pkg/concurrency"
	"github.com/lxdlam/vertex/pkg/protocol"
)

// Request warps the request redis object and the target id
type Request struct {
	ClientID int
	Req      protocol.RedisObject
}

// Response warps a redis object and the target id
type Response struct {
	ClientID int
	Resp     protocol.RedisObject
}

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
	eventBus         concurrency.EventBus
	shutChan         chan os.Signal
	shutdown         int32
	responseReceiver concurrency.Receiver
	cleanUpHandles   []func()
	clients          sync.Map
}

//func NewServer() Server {
//}

func (s *server) Init(c common.Config) bool {
	var err error

	common.Debug("Initialize server")
	defer common.Debug("Initialize server success")

	addr := &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: c.Port,
	}

	s.tcpListener, err = net.ListenTCP("tcp", addr)
	if err != nil {

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

	common.Infof("Server is listening %s", addr.String())

	return true
}

func (s *server) Serve() {
	var wg sync.WaitGroup
	wg.Add(1)

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
				s.newConn(conn)
			}
		}

		wg.Done()
	}()

	s.clean()
	// wait until shutdown
	wg.Wait()
}

func (s *server) Stop() {
	if atomic.CompareAndSwapInt32(&s.shutdown, 0, 1) {
		s.shutChan <- syscall.SIGTERM
	}
}

func (s *server) newConn(conn net.Conn) {

}

func (s *server) response(resp Response) {
	conn, ok := s.clients.Load(resp.ClientID)
	if !ok {
		common.Errorf("client %d not exist", resp.ClientID)
	}

	respByte := resp.Resp.Byte()

	write, err := conn.(*net.TCPConn).Write(respByte)

}

func (s *server) clean() {

}
