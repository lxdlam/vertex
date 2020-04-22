package concurrency

// BroadCaster is a simple message wrapper.
//
// It provides the semantic for 1 to M message consumers, no error will be produced.
// It's goroutine safe but a one shot message producer which means multiple set will make the value unreliable.
type BroadCaster interface {
	// Set the message, set a nil value also makes sense.
	Set(interface{})

	// Get will blocking wait for a new message or return the received message
	Get() interface{}
}

type broadCaster struct {
	message      interface{}
	readyChannel chan bool
}

// NewBroadCaster will return a BroadCaster instance.
func NewBroadCaster() BroadCaster {
	return &broadCaster{
		message:      nil,
		readyChannel: make(chan bool),
	}
}

func (bc *broadCaster) Set(message interface{}) {
	bc.message = message
	close(bc.readyChannel)
}

func (bc *broadCaster) Get() interface{} {
	<-bc.readyChannel
	return bc.message
}
