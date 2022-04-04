package pubsub

import (
	"sync"
)

type wrappedMessage[MessageType interface{}] struct {
	topic   string
	message MessageType
}

type Engine[MessageType interface{}] struct {
	messages chan *wrappedMessage[MessageType]

	subscribersByTopicLock *sync.Mutex
	subscribersByTopic     map[string][]chan MessageType
}

func NewEngine[MessageType interface{}](capacity int) *Engine[MessageType] {
	return &Engine[MessageType]{
		messages: make(chan *wrappedMessage[MessageType], capacity),

		subscribersByTopicLock: &sync.Mutex{},
		subscribersByTopic:     make(map[string][]chan MessageType),
	}
}

func (e *Engine[MessageType]) Start() {
	for wm := range e.messages {
		e.broadcast(wm.topic, wm.message)
	}
}

func (e *Engine[MessageType]) Stop() {
	close(e.messages)
}

func (e *Engine[MessageType]) Publish(topic string, msg MessageType) {
	e.messages <- &wrappedMessage[MessageType]{
		topic:   topic,
		message: msg,
	}
}

func (e *Engine[MessageType]) Subscribe(topic string, capacity int) <-chan MessageType {
	e.subscribersByTopicLock.Lock()
	defer e.subscribersByTopicLock.Unlock()

	subscribers, has := e.subscribersByTopic[topic]
	if !has {
		subscribers = make([]chan MessageType, 0)
		e.subscribersByTopic[topic] = subscribers
	}

	s := make(chan MessageType, capacity)
	subscribers = append(subscribers, s)
	e.subscribersByTopic[topic] = subscribers

	return s
}

func (e *Engine[MessageType]) Unsubscribe(topic string, s <-chan MessageType) {
	e.subscribersByTopicLock.Lock()
	defer e.subscribersByTopicLock.Unlock()

	subscribers, has := e.subscribersByTopic[topic]
	if !has {
		return
	}

	indexToRemove := -1

	for i, sub := range subscribers {
		if sub == s {
			indexToRemove = i
		}
	}

	if indexToRemove < 0 {
		return
	}

	subscribers = append(subscribers[:indexToRemove], subscribers[indexToRemove+1:]...)

	e.subscribersByTopic[topic] = subscribers
}

func (e *Engine[MessageType]) broadcast(topic string, m MessageType) {
	e.subscribersByTopicLock.Lock()
	defer e.subscribersByTopicLock.Unlock()

	subscribers, ok := e.subscribersByTopic[topic]
	if !ok {
		return
	}

	for _, s := range subscribers {
		s <- m
	}
}
