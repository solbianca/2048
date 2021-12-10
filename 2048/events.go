package twenty48

// An Event Bus is an implementation of the pub/sub pattern where publishers are publishing data and
// interested subscribers can listen to them and act based on data.
// Based on article https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997

import (
	"sync"
)

var Events *EventBus

func init() {
	Events = NewEventBus()
}

type DataEvent struct {
	Data  interface{}
	Topic string
}

// DataChannel is a channel which can accept an DataEvent
type DataChannel chan DataEvent

// DataChannelSlice is a slice of DataChannels
type DataChannelSlice []DataChannel

// EventBus stores the information about subscribers interested for a particular topic
type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{subscribers: map[string]DataChannelSlice{}}
}

func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()
	if chans, found := eb.subscribers[topic]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		// special thanks for /u/freesid who pointed it out
		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
	eb.rm.RUnlock()
}

func (eb *EventBus) Subscribe(topic string, ch DataChannel) {
	eb.rm.Lock()
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()
}

func Publish(topic string, data interface{}) {
	Events.Publish(topic, data)
}

func Subscribe(topic string, ch DataChannel) {
	Events.Subscribe(topic, ch)
}
