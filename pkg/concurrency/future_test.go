package concurrency_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/lxdlam/vertex/pkg/concurrency"
)

var (
	testError    = errors.New("TestError")
	shortTime    = 100 * time.Millisecond
	longTime     = 1 * time.Second
	maxGoroutine = 10
)

// The only test for WithoutErrorFuture
func TestNoErrorFutureGet(t *testing.T) {
	fut := NewFuture(NewNoErrorTask(func() interface{} {
		time.Sleep(shortTime)
		return 1
	}))

	if !testGet(t, fut, 1, nil) {
		t.Fatalf("Fail at first get.")
	}

	// Multiple get should be same
	if !testGet(t, fut, 1, nil) {
		t.Fatalf("Fail at second get.")
	}
}

// The only test for ErrorFuture
func TestErrorFutureGet(t *testing.T) {
	fut := NewFuture(NewErrorTask(testError))

	if !testGet(t, fut, nil, testError) {
		t.Fatalf("Fail at first get.")
	}

	// Multiple get should be same
	if !testGet(t, fut, nil, testError) {
		t.Fatalf("Fail at second get.")
	}
}

func TestFutureGetNoError(t *testing.T) {
	fut := newFuture(1, nil, shortTime)

	if !testGet(t, fut, 1, nil) {
		t.Fatalf("Fail at first get.")
	}

	// Multiple get should be same
	if !testGet(t, fut, 1, nil) {
		t.Fatalf("Fail at second get.")
	}
}

func TestFutureGetWithError(t *testing.T) {
	fut := newFuture(nil, testError, shortTime)

	if !testGet(t, fut, nil, testError) {
		t.Fatalf("Fail at first get.")
	}

	// Multiple get should be same
	if !testGet(t, fut, nil, testError) {
		t.Fatalf("Fail at second get.")
	}
}

func TestFutureCancel(t *testing.T) {
	fut := newFuture(1, nil, shortTime)

	err := fut.Cancel()
	if !assert.Nil(t, err) {
		t.Fatalf("The first cancel should return no error, got=%+v", err)
	}

	err = fut.Cancel()
	if !assert.Equal(t, ErrCancelled, err) {
		t.Fatalf("The second cancel should return ErrCancelled, got=%+v", err)
	}

	// Should get with nil and ErrCancelled
	if !testGet(t, fut, nil, ErrCancelled) {
		t.Fatalf("Get a cancelled future with unexpected result.")
	}

	err = fut.Wait()
	if !assert.Equal(t, ErrCancelled, err) {
		t.Fatalf("The wait call should return ErrCancelled, got=%+v", err)
	}

	err = fut.WaitFor(shortTime)
	if !assert.Equal(t, ErrCancelled, err) {
		t.Fatalf("The wait for call should return ErrCancelled, got=%+v", err)
	}

	fut = NewFuture(NewTask(func() (interface{}, error) {
		time.Sleep(shortTime)

		return 1, nil
	}))

	_ = fut.Wait()

	err = fut.Cancel()
	if !assert.Equal(t, ErrFulfilled, err) {
		t.Fatalf("Cancel a fullfilled future should return ErrCancelled, got=%+v", err)
	}

	if !testGet(t, fut, 1, nil) {
		t.Fatalf("Get from a fullfilled future failed.")
	}
}

func TestCancelAndFullfilConflict(t *testing.T) {
	// A very special case
	// If we just cancelled a future, then the future is finished
	// Then we cancel again, what error should be reported? ErrCancelled.

	var wg sync.WaitGroup

	// The first case, cancel it ASAP, then check if it's reporting the ErrCancelled
	wg.Add(1)

	fut := NewFuture(NewTask(func() (interface{}, error) {
		time.Sleep(shortTime)

		defer wg.Done()
		return 1, nil
	}))

	err := fut.Cancel()
	if !assert.Nil(t, err) {
		t.Fatalf("Cancel a running future should return no error, got=%+v", err)
	}

	wg.Wait()
	time.Sleep(time.Millisecond)

	for i := 0; i < 10; i++ {
		err = fut.Cancel()
		if !assert.Equal(t, ErrCancelled, err) {
			t.Fatalf("Cancel a already cancelled future should return ErrCancelled, got=%+v", err)
		}
	}

	// Then the second case, just waiting for the task finish, then cancel should return ErrFullfilled
	fut = NewFuture(NewTask(func() (interface{}, error) {
		time.Sleep(shortTime)

		return 1, nil
	}))

	_ = fut.Wait()

	for i := 0; i < 10; i++ {
		err = fut.Cancel()
		if !assert.Equal(t, ErrFulfilled, err) {
			t.Fatalf("Cancel a fullfilled future should return ErrCancelled, got=%+v", err)
		}
	}
}

func TestWaitFor(t *testing.T) {
	fut := newFuture(1, nil, longTime)

	err := fut.WaitFor(shortTime)
	if !assert.Equal(t, ErrTimeout, err) {
		t.Fatalf("WaitFor reaches time limit should return ErrTimeout, got=%+v", err)
	}

	_ = fut.Wait()

	err = fut.WaitFor(shortTime)
	if !assert.Nil(t, err) {
		t.Fatalf("WaitFor wait for a fullfilled future should return no error, got=%+v", err)
	}
}

func TestGetMultipleGoroutine(t *testing.T) {
	fut := newFuture(1, nil, longTime)

	for i := 0; i < maxGoroutine; i++ {
		go func() {
			if !testGet(t, fut, 1, nil) {
				t.Fatalf("Fail at get.")
			}
		}()
	}
}

func TestCancelMultipleGoroutine(t *testing.T) {
	fut := newFuture(1, nil, longTime)

	for i := 0; i < maxGoroutine; i++ {
		go func() {
			if !testGet(t, fut, nil, ErrCancelled) {
				t.Fatalf("Fail at get.")
			}
		}()

		go func() {
			err := fut.Wait()
			if !assert.Equal(t, ErrCancelled, err) {
				t.Fatalf("The wait call should return ErrCancelled, got=%+v", err)
			}
		}()

		go func() {
			err := fut.WaitFor(longTime)
			if !assert.Equal(t, ErrCancelled, err) {
				t.Fatalf("The wait call should return ErrCancelled, got=%+v", err)
			}
		}()
	}

	time.Sleep(shortTime)

	err := fut.Cancel()
	if !assert.Nil(t, err) {
		t.Fatalf("The first cancel should return no error, got=%+v", err)
	}

	err = fut.Cancel()
	if !assert.Equal(t, ErrCancelled, err) {
		t.Fatalf("The second cancel should return ErrCancelled, got=%+v", err)
	}
}

func TestWaitForMultipleGoroutine(t *testing.T) {
	fut := newFuture(1, nil, longTime*2)

	for i := 0; i < maxGoroutine; i++ {
		go func() {
			err := fut.WaitFor(shortTime)
			if !assert.Equal(t, ErrTimeout, err) {
				t.Fatalf("WaitFor reaches time limit should return ErrTimeout, got=%+v", err)
			}
		}()

		go func() {
			err := fut.WaitFor(longTime)
			if !assert.Equal(t, ErrTimeout, err) {
				t.Fatalf("WaitFor reaches time limit should return ErrTimeout, got=%+v", err)
			}
		}()
	}

	_ = fut.Wait()

	for i := 0; i < maxGoroutine; i++ {
		go func() {
			err := fut.WaitFor(shortTime)
			if !assert.Nil(t, err) {
				t.Fatalf("WaitFor wait for a fullfilled future should return no error, got=%+v", err)
			}
		}()

		go func() {
			err := fut.WaitFor(shortTime)
			if !assert.Nil(t, err) {
				t.Fatalf("WaitFor wait for a fullfilled future should return no error, got=%+v", err)
			}
		}()
	}
}

func newFuture(value interface{}, err error, d time.Duration) Future {
	return NewFuture(NewTask(func() (interface{}, error) {
		time.Sleep(d)
		return value, err
	}))
}

func testGet(t *testing.T, fut Future, expected_value interface{}, expected_error error) bool {
	actual_value, actual_error := fut.Get()

	if !assert.Equal(t, expected_value, actual_value, "The given value is not equal.") {
		return false
	}
	if !assert.Equal(t, expected_error, actual_error, "The given error is not equal.") {
		return false
	}

	return true
}
