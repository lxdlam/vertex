package concurrency

import (
	"errors"
	"sync"
	"time"

	"github.com/lxdlam/vertex/pkg/util"
)

const (
	anonymousSource string = "anonymous"
)

var (
	// ErrNoSuchTopic will be raised if no topic is found of the given name
	ErrNoSuchTopic = errors.New("event_bus: no such topic")
)

// Event includes all necessary data to describe a Event.
type Event interface {
	// Data will return the payload of the event
	Data() interface{}

	// ID will return the UUID of the event
	ID() string

	// Topic will return the source topic of the Event
	Topic() string

	// Time will return the Unix timestamp of when the event has occurred.
	Time() int64

	// Source is a additional property to report the sender of the event
	// It is empty by default.
	Source() string
}

type event struct {
	id     string
	data   interface{}
	topic  string
	time   int64
	source string
}

func (e *event) Data() interface{} {
	return e.data
}

func (e *event) ID() string {
	return e.id
}

func (e *event) Topic() string {
	return e.topic
}

func (e *event) Time() int64 {
	return e.time
}

func (e *event) Source() string {
	return e.source
}

// NewEvent will return a new event instance, set source as anonymous
func NewEvent(topic string, data interface{}) Event {
	return NewEventWithSource(topic, anonymousSource, data)
}

// NewEventWithSource will return a new event instance with the given source
func NewEventWithSource(topic, source string, data interface{}) Event {
	return &event{
		id:     util.GenNewUUID(),
		data:   data,
		topic:  topic,
		time:   time.Now().Unix(),
		source: source,
	}
}

// EventBus is a channel based, goroutine safe event distribute util
type EventBus interface {
	// Publish will send a Event to the specific topic, returns a future immediately:
	// - If success, the future will return the number of successfully delivered receivers
	// - Otherwise the err will be set.
	// The source of Publish will be "anonymous"
	Publish(string, interface{}) Future

	// PublishWithSource will publish a event with a specific source. It works same as Publish.
	PublishWithSource(string, string, interface{}) Future

	// Subscribe will get a receiver of the corresponding topic.
	// If no topic with the given name, the error will be ErrNoSuchTopic.
	// No unsubscribe semantic will be given, the receiver can self unsubscribe.
	Subscribe(string, string) (Receiver, error)

	// NewTopic will generate a new topic, false if topic exists
	NewTopic(string) bool

	// RemoveTopic will remove the given topic, false if no such topic exists
	RemoveTopic(string) bool

	// ExistTopic will just simply check if the topic is exist
	ExistTopic(string) bool
}

type eventBus struct {
	// Design note: topic will not been modified so frequently, so just using a simple RWMutex to protect it
	mutex  sync.RWMutex
	topics map[string]Topic
}

func (e *eventBus) Publish(name string, data interface{}) Future {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	topic, ok := e.topics[name]

	if !ok {
		return NewFuture(NewErrorTask(ErrNoSuchTopic))
	}

	return topic.Publish(NewEvent(name, data))
}

func (e *eventBus) PublishWithSource(name, source string, data interface{}) Future {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	topic, ok := e.topics[name]

	if !ok {
		return NewFuture(NewErrorTask(ErrNoSuchTopic))
	}

	return topic.Publish(NewEventWithSource(name, source, data))
}

func (e *eventBus) Subscribe(name string, subscriber string) (Receiver, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	topic, ok := e.topics[name]

	if !ok {
		return nil, ErrNoSuchTopic
	}

	return topic.Subscribe(subscriber), nil

}

func (e *eventBus) NewTopic(name string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	_, exist := e.topics[name]

	if exist {
		return false
	}

	e.topics[name] = NewTopic(name)

	return true
}

func (e *eventBus) RemoveTopic(name string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	topic, exist := e.topics[name]

	if !exist {
		return false
	}

	topic.Remove()
	delete(e.topics, name)

	return true
}

func (e *eventBus) ExistTopic(name string) bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	_, exist := e.topics[name]

	return exist
}
