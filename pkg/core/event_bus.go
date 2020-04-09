package core

import (
	"sync"
)

// Event includes all necessary for EventBus to transfer
type Event struct {
	ID        string
	Payload   interface{}
	TopicName string
	Time      int64
}

type topic struct {
}

// EventBus is an internal communicate medium across all components
type EventBus struct {
	mutex sync.Mutex
}
