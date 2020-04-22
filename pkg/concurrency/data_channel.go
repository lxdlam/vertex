package concurrency

import (
	"errors"
	"sync/atomic"
	"time"
)

var (
	// ErrChannelClosed will be raised if we access after a channel is closed
	ErrChannelClosed = errors.New("data channel: the channel has already closed")
)

const (
	// DefaultExpireTime is the default time limit for send operation
	DefaultExpireTime = 10 * time.Millisecond
)

const (
	// Success means the event is successfully delivered
	Success int = iota

	// Closed means the receiver has closed the channel
	Closed

	// Expired means the sender waits too long and reaches the time limit,
	// then the event will be dropped.
	Expired
)

// Receiver is a receiver side of the data channel.
//
// Receiver can safely close the channel.
type Receiver interface {
	// Receive works like <-ch, will blocking until the data is received.
	// If the channel is closed, ErrChannelClosed will be returned.
	Receive() (Event, error)

	// Close the channel.
	// For buffered data channel, all the unread data will be dropped after close which
	// means after close, you will always get ErrChannelClosed regardless what type of data channel
	// you are using.
	Close()
}

// Sender is the send side of the data channel.
//
// No close semantic is provided from this side.
type Sender interface {
	// Send a event to the corresponding receiver.
	// - Closed if you're sending to a closed channel
	// - Expired if you reached the time limit
	// - Otherwise success
	Send(Event) int
}

type receiver struct {
	dataChannel  <-chan Event
	closeChannel chan<- struct{}
	closed       int32
}

type sender struct {
	dataChannel  chan<- Event
	closeChannel <-chan struct{}
	expireTime   time.Duration
}

// NewDataChannel will generate a receiver and a sender pair.
func NewDataChannel() (Receiver, Sender) {
	dc := make(chan Event)
	cc := make(chan struct{})

	r := &receiver{
		dataChannel:  dc,
		closeChannel: cc,
		closed:       0,
	}

	s := &sender{
		dataChannel:  dc,
		closeChannel: cc,
		expireTime:   DefaultExpireTime,
	}

	return r, s
}

// NewDataChannelWithOption will return a DataChannel with the given options
// If the size is less or equal 0, the buffer size will be infinity. You can use DefaultExpireTime to set the second one.
func NewDataChannelWithOption(size int, expireTime time.Duration) (Receiver, Sender) {
	var dc chan Event
	if size <= 0 {
		dc = make(chan Event)
	} else {
		dc = make(chan Event, size)
	}
	cc := make(chan struct{})

	r := &receiver{
		dataChannel:  dc,
		closeChannel: cc,
		closed:       0,
	}

	s := &sender{
		dataChannel:  dc,
		closeChannel: cc,
		expireTime:   expireTime,
	}

	return r, s
}

func (r *receiver) Receive() (Event, error) {
	if atomic.LoadInt32(&r.closed) == 1 {
		return nil, ErrChannelClosed
	}

	e := <-r.dataChannel

	return e, nil
}

func (r *receiver) Close() {
	if atomic.CompareAndSwapInt32(&r.closed, 0, 1) {
		close(r.closeChannel)
	}
}

func (s *sender) Send(event Event) int {
	select {
	case <-s.closeChannel:
		return Closed
	case s.dataChannel <- event:
		return Success
	case <-time.After(s.expireTime):
		return Expired
	}
}
