package core

// Server is the main service of vertex.
type Server interface {
	// Prepare the resource, false if failed.
	// If it returns false, the server should not be start.
	Init() bool

	// Serve call will start the run session with blocking the current goroutine.
	// If it returns, it means that the server is requested to be stopped.
	Serve()

	// Stop will stop the server. It provides a way to manually stop the server.
	Stop()
}

//func NewServer() Server {
//}

type server struct {
	cleanUpHandles []func()
}

func (s *server) Init() bool {
	panic("implement me")
}

func (s *server) Serve() {
	panic("implement me")
}

func (s *server) Stop() {
	panic("implement me")
}
