package core

import (
	"errors"
	"github.com/lxdlam/vertex/pkg/common"
	"github.com/lxdlam/vertex/pkg/concurrency"
	"sync"
	"sync/atomic"
)

var (
	// ErrTopicRemoved will be sent to all receiver when a topic is removed from data bus
	ErrTopicRemoved = errors.New("topic: topic has been removed")
)

// Topic
type Topic interface {
	// Subscriber will return a new Receiver and register the Sender to itself.
	Subscribe(string) Receiver

	// Publish will distribute all Events to the registered Sender.
	// When distributing, it will start goroutines and wait them to expire.
	Publish(Event) concurrency.Future

	// Remove will do tht send work when the removing the topic from data bus.
	Remove()
}

func NewTopic(name string) Topic {
	return &topic{
		name: name,
	}
}

type topic struct {
	name        string
	subscribers sync.Map
}

func (t *topic) Subscribe(name string) Receiver {
	r, s := NewDataChannelWithOption(10, DefaultExpireTime)
	t.subscribers.Store(name, s)
	return r
}

func (t *topic) Publish(e Event) (fut concurrency.Future) {
	ch := make(chan int32, 1)

	fut = concurrency.NewFuture(concurrency.NewTask(func() (interface{}, error) {
		return <-ch, nil
	}))

	go t.distributeEvent(e, ch)

	return
}

func (t *topic) Remove() {
	ch := make(chan int32, 1)

	go t.distributeEvent(NewEvent(t.name, ErrTopicRemoved), ch)

	<-ch
}

func (t *topic) distributeEvent(e Event, successChan chan<- int32) {
	var closedKeys []string
	var wg sync.WaitGroup
	var successCount int32

	ch := make(chan string)

	go func() {
		for key := range ch {
			closedKeys = append(closedKeys, key)
		}
	}()

	t.subscribers.Range(func(k interface{}, s interface{}) bool {
		wg.Add(1)

		name := k.(string)
		subscriber := s.(Sender)

		go func() {
			status := subscriber.Send(e)
			if status == Closed {
				ch <- name
			} else if status == Expired {
				common.Infof("send event to subscriber expired. topic=%s, subscriber=%s", t.name, name)
			} else {
				atomic.AddInt32(&successCount, 1)
			}

			wg.Done()
		}()

		return true
	})

	wg.Wait()
	successChan <- successCount
	close(ch)

	go t.batchRemove(closedKeys)
}

func (t *topic) batchRemove(removeKeys []string) {
	for _, name := range removeKeys {
		t.subscribers.Delete(name)
	}
}
