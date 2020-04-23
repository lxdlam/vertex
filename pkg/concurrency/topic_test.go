package concurrency_test

import (
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/lxdlam/vertex/pkg/concurrency"
)

const (
	subscribersCount int = 5
)

func TestSubscribeAndPublish(t *testing.T) {
	var wg sync.WaitGroup

	actualSlice := make([]int64Slice, subscribersCount)
	subscribers := make([]Receiver, subscribersCount)
	topic := NewTopic("test")
	data := genTestData(100)

	wg.Add(1 + subscribersCount)

	for i := 0; i < subscribersCount; i++ {
		subscribers[i] = topic.Subscribe(fmt.Sprintf("subscriber%d", i))
		go func(idx int) {
			for i := 0; i < 100; i++ {
				event, err := subscribers[idx].Receive()
				assert.Nil(t, err)
				assert.Nil(t, event.Error())

				actualSlice[idx] = append(actualSlice[idx], event.Data().(int64))
			}

			subscribers[idx].Close()
			wg.Done()
		}(i)
	}

	go func() {
		var futures []Future
		for _, num := range data {
			futures = append(futures, topic.Publish(newEvent(num)))
		}

		for _, future := range futures {
			count, err := future.Get()
			assert.Equal(t, subscribersCount, count)
			assert.Nil(t, err)
		}

		wg.Done()
	}()

	wg.Wait()

	expectedSlice := int64Slice(data)

	sort.Sort(expectedSlice)

	for i := 0; i < subscribersCount; i++ {
		sort.Sort(actualSlice[i])
		assert.ElementsMatch(t, expectedSlice, actualSlice[i])
	}
}

func TestSubscriberLeaveAndClean(t *testing.T) {
	var stage1, stage2, internal, all sync.WaitGroup

	topic := NewTopic("test")
	subscriber1 := topic.Subscribe("s1")
	subscriber2 := topic.Subscribe("s2")

	all.Add(3)
	stage1.Add(2)
	stage2.Add(1)
	internal.Add(1)

	go func() {
		event, err := subscriber1.Receive()
		assert.Nil(t, err)
		assert.Equal(t, 1, event.Data().(int))
		assert.Nil(t, event.Error())

		subscriber1.Close()
		stage1.Done()
		internal.Done()

		event, err = subscriber1.Receive()
		assert.Equal(t, ErrChannelClosed, err)
		assert.Nil(t, event)

		all.Done()
	}()

	go func() {
		event, err := subscriber2.Receive()
		assert.Nil(t, err)
		assert.Equal(t, 1, event.Data().(int))
		assert.Nil(t, event.Error())

		stage1.Done()
		internal.Wait()

		event, err = subscriber2.Receive()
		assert.Nil(t, err)
		assert.Equal(t, 2, event.Data().(int))
		assert.Nil(t, event.Error())

		subscriber2.Close()
		stage2.Done()

		event, err = subscriber2.Receive()
		assert.Equal(t, ErrChannelClosed, err)
		assert.Nil(t, event)

		all.Done()
	}()

	go func() {
		fut := topic.Publish(newEvent(1))
		ret, err := fut.Get()
		assert.Equal(t, 2, ret)
		assert.Nil(t, err)

		stage1.Wait()

		fut = topic.Publish(newEvent(2))
		ret, err = fut.Get()
		assert.Equal(t, 1, ret)
		assert.Nil(t, err)

		stage2.Wait()

		fut = topic.Publish(newEvent(3))
		ret, err = fut.Get()
		assert.Equal(t, 0, ret)
		assert.Nil(t, err)

		all.Done()
	}()

	all.Wait()
}

func TestTopicRemove(t *testing.T) {
	var wg sync.WaitGroup

	topic := NewTopic("test")
	subscribers := make([]Receiver, subscribersCount)

	wg.Add(1 + subscribersCount)

	for i := 0; i < subscribersCount; i++ {
		subscribers[i] = topic.Subscribe(fmt.Sprintf("subscriber%d", i))
		go func(idx int) {
			event, err := subscribers[idx].Receive()
			assert.Equal(t, ErrTopicRemoved, event.Error().(error))
			assert.Nil(t, event.Data())
			assert.Nil(t, err)

			wg.Done()
		}(i)
	}

	go func() {
		topic.Remove()
		wg.Done()
	}()

	wg.Wait()
}
