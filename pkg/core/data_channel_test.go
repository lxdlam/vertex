package core_test

import (
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/lxdlam/vertex/pkg/core"
)

type int64Slice []int64

func (is int64Slice) Len() int {
	return len(is)
}

func (is int64Slice) Less(i, j int) bool {
	return is[i] < is[j]
}

func (is int64Slice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

func newEvent(data interface{}) Event {
	return NewEvent("dummy", data)
}

func TestDataChannelConcurrentSendReceive(t *testing.T) {
	r, s := NewDataChannel()
	var expectedSlice, actualSlice int64Slice
	var wg sync.WaitGroup
	wg.Add(2)

	// Receive side
	go func() {
		for {
			if val, err := r.Receive(); err != nil {
				t.Fatalf("receive failed, err=%v", err)
			} else {
				actualSlice = append(actualSlice, val.Data().(int64))
			}

			if len(actualSlice) == 10000 {
				break
			}
		}

		wg.Done()
	}()

	// Send side
	go func() {
		vis := make(map[int64]bool)
		rand.Seed(time.Now().UnixNano())

		for i := 0; i < 10000; i++ {
			var num int64
			for {
				num = rand.Int63()
				if _, ok := vis[num]; !ok {
					vis[num] = true
					break
				}
			}

			expectedSlice = append(expectedSlice, num)

			go func(sl int, n int64) {
				time.Sleep(time.Duration(sl) * time.Millisecond)
				s.Send(newEvent(n))
			}(rand.Intn(100)+1, num)
		}

		wg.Done()
	}()

	wg.Wait()

	sort.Sort(expectedSlice)
	sort.Sort(actualSlice)

	assert.ElementsMatch(t, expectedSlice, actualSlice)
}

func TestDataChannelSendAndReceiveClosed(t *testing.T) {
	r, s := NewDataChannel()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		time.Sleep(10 * time.Millisecond)
		r.Close()
		wg.Done()
	}()

	go func() {
		// Send will blocking here, and the receiver will close earlier
		status := s.Send(newEvent(15))
		assert.Equal(t, Closed, status)
		wg.Done()
	}()

	wg.Wait()

	_, err := r.Receive()
	assert.Equal(t, err, ErrChannelClosed)
}

func TestDataBufferedChannelSendClosed(t *testing.T) {
	r, s := NewDataChannelWithOption(10, DefaultExpireTime)
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan struct{})

	go func() {
		time.Sleep(150 * time.Millisecond)
		r.Close()
		close(ch)
		wg.Done()
	}()

	go func() {
		counter := 0
	outer:
		for {
			select {
			case <-ch:
				break outer
			case <-time.After(10 * time.Millisecond):
				status := s.Send(newEvent(15))
				if counter < 10 {
					// In a buffered data channel, the first 10 event will be sent immediately
					counter += 1
					assert.Equal(t, Success, status)
				} else {
					// But the 11th will fail
					assert.Equal(t, Closed, status)
				}
			}
		}

		wg.Done()
	}()

	wg.Wait()

	// The closed channel will produce ErrChannelClosed regardless of type
	_, err := r.Receive()
	assert.Equal(t, err, ErrChannelClosed)
}

func TestMultipleSendCloseSafety(t *testing.T) {
	r, s := NewDataChannel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		s.Send(newEvent(1))
		wg.Done()
	}()

	go func() {
		s.Send(newEvent(2))
		wg.Done()
	}()

	// No panic should produced
	r.Close()
	wg.Wait()
}

func TestExpiredSend(t *testing.T) {
	r, s := NewDataChannelWithOption(0, 10*time.Microsecond)

	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan struct{})

	go func() {
		status := s.Send(newEvent(15))
		assert.Equal(t, status, Expired)
		<-ch
		status = s.Send(newEvent(15))
		assert.Equal(t, status, Closed)

		wg.Done()
	}()

	time.Sleep(20 * time.Millisecond)

	r.Close()
	close(ch)
	wg.Wait()
}
