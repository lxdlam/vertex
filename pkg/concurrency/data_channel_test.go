package concurrency_test

import (
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/lxdlam/vertex/pkg/concurrency"
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
	return NewEvent("dummy", data, nil)
}

func genTestData(size int) []int64 {
	rand.Seed(time.Now().UnixNano())

	var result []int64
	vis := make(map[int64]bool)
	for i := 0; i < size; i++ {
		var num int64
		for {
			num = rand.Int63()
			if _, ok := vis[num]; !ok {
				vis[num] = true
				result = append(result, num)
				break
			}
		}
	}

	return result
}

func TestDataChannelConcurrentSendReceive(t *testing.T) {
	r, s := NewDataChannel()
	var expectedSlice, actualSlice int64Slice
	var wg sync.WaitGroup

	data := genTestData(100)
	wg.Add(2)

	// Receive side
	go func() {
		for i := 0; i < 100; i++ {
			if val, err := r.Receive(); err != nil {
				t.Fatalf("receive failed, err=%v", err)
			} else {
				actualSlice = append(actualSlice, val.Data().(int64))
			}
		}

		wg.Done()
	}()

	// Send side
	go func() {
		for _, num := range data {
			go func(sl int, n int64) {
				time.Sleep(time.Duration(sl) * time.Microsecond)
				status := s.Send(newEvent(n))
				assert.Equal(t, Success, status)
			}(rand.Intn(10)+1, num)

			expectedSlice = append(expectedSlice, num)
		}

		wg.Done()
	}()

	wg.Wait()

	sort.Sort(expectedSlice)
	sort.Sort(actualSlice)

	assert.ElementsMatch(t, expectedSlice, actualSlice)
}

func TestDataChannelSendAndReceiveClosed(t *testing.T) {
	r, s := NewDataChannelWithOption(0, time.Second)
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
	r, s := NewDataChannelWithOption(10, time.Second)
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan struct{})

	go func() {
		<-ch
		r.Close()
		wg.Done()
	}()

	go func() {
		for i := 0; i < 10; i++ {
			status := s.Send(newEvent(15))
			assert.Equal(t, Success, status)
		}

		close(ch)

		status := s.Send(newEvent(15))
		assert.Equal(t, Closed, status)

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
	r, s := NewDataChannelWithOption(0, 1*time.Millisecond)

	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan struct{})

	go func() {
		status := s.Send(newEvent(15))
		assert.Equal(t, Expired, status)
		<-ch
		status = s.Send(newEvent(15))
		assert.Equal(t, Closed, status)

		wg.Done()
	}()

	time.Sleep(20 * time.Millisecond)

	r.Close()
	close(ch)
	wg.Wait()
}
