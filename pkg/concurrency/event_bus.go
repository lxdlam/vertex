package concurrency

import (
	"sync"
)

// Event includes all necessary for EventBus to transfer
type Event struct {
	ID      string
	Payload interface{}
	Topic   string
	Time    int64
	From    string
}

type topic struct {
	mutex       sync.RWMutex
	name        string
	subscribers []string
}

func newTopic(name string) topic {
	return topic{
		name: name,
	}
}

// EventBus is a channel based, event distribute util.
//
// The event bus is goroutine safe, so does the topic.
type EventBus interface {
	Publish(string, Event) bool
	AsyncPublish(string, Event) Future
	Subscribe(string) bool
	Unsubscribe(string) bool
}

type eventBus struct {
	mutex        sync.RWMutex
	shutdownChan chan bool
	topics       map[string]topic
}
