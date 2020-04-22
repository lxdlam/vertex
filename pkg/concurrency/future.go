package concurrency

import (
	"errors"
	"time"
)

var (
	// ErrCancelled will happen when invoke access functions on a cancelled future
	ErrCancelled = errors.New("future: The future has already been cancelled")

	// ErrFulfilled will happen when invoke cancel function on a fulfilled future
	ErrFulfilled = errors.New("future: The future has already fulfilled")

	// ErrTimeout will happen when waiting a future reaches the timeout
	ErrTimeout = errors.New("future: waiting time out")
)

// Future is a simple goroutine wrapper for some async works.
//
// The terminology Future is same as futures in C++ or Rust, provides cancel and wait operations to control the task.
type Future interface {
	// Cancel a Future, error will be raised if:
	// - The future has already been cancelled
	// - The future has already been fulfilled
	Cancel() error

	// Get the value of the future mission, will blocking until the task is finished.
	// If the future has been cancelled, the result will be nil and error will be set.
	// Otherwise the result and the error raised by the task will return.
	Get() (interface{}, error)

	// Wait until the future will be fulfilled.
	// If the future has already fulfilled, the wait will return immediately.
	// If the Cancel() has been invoked when waiting, an error will be returned.
	Wait() error

	// WaitFor works like the wait but will timeout when reach the given time duration.
	// A special error will be returned when timeout.
	WaitFor(d time.Duration) error
}

type future struct {
	task       Task
	result     interface{}
	err        error
	doneChan   chan struct{}
	cancelChan chan struct{}
}

// NewFuture will start a new future with the given function.
// The function will be invoked immediately when successfully created.
func NewFuture(task Task) Future {
	fut := &future{
		task:       task,
		result:     nil,
		err:        nil,
		doneChan:   make(chan struct{}),
		cancelChan: make(chan struct{}),
	}

	fut.start()

	return fut
}

func (f *future) start() {
	go func() {
		f.result, f.err = f.task.Run()
		select {
		// Already cancelled, leaves the doneChan open
		case <-f.cancelChan:
			return
		default:
			close(f.doneChan)
		}
	}()
}

func (f *future) Cancel() error {
	select {
	case <-f.doneChan:
		return ErrFulfilled
	case <-f.cancelChan:
		return ErrCancelled
	default:
		close(f.cancelChan)
		return nil
	}
}

func (f *future) Get() (interface{}, error) {
	err := f.Wait()

	if err != nil {
		return nil, err
	}

	return f.result, f.err
}

func (f *future) Wait() error {
	select {
	case <-f.doneChan:
		return nil
	case <-f.cancelChan:
		return ErrCancelled
	}
}

func (f *future) WaitFor(d time.Duration) error {
	select {
	case <-f.doneChan:
		return nil
	case <-f.cancelChan:
		return ErrCancelled
	case <-time.After(d):
		return ErrTimeout
	}
}
