package pubsub

import (
	"sync"
)

type Engine struct {
	messages chan *wrappedMessage

	subscribersByTopicLock *sync.Mutex
	subscribersByTopic     map[string][]chan *Message
}

func NewEngine(capacity int) *Engine {
	return &Engine{
		messages: make(chan *wrappedMessage, capacity),

		subscribersByTopicLock: &sync.Mutex{},
		subscribersByTopic:     make(map[string][]chan *Message),
	}
}

func (e *Engine) Start() {
	for wm := range e.messages {
		e.broadcast(wm.topic, wm.message)
	}
}

func (e *Engine) Stop() {
	close(e.messages)
}

func (e *Engine) Publish(topic string, name string, data interface{}) {
	m := &Message{
		Name: name,
		Data: data,
	}

	e.PublishMessage(topic, m)
}

func (e *Engine) PublishMessage(topic string, m *Message) {
	e.messages <- &wrappedMessage{
		topic:   topic,
		message: m,
	}
}

func (e *Engine) Subscribe(topic string, capacity int) <-chan *Message {
	e.subscribersByTopicLock.Lock()
	defer e.subscribersByTopicLock.Unlock()

	subscribers, has := e.subscribersByTopic[topic]
	if !has {
		subscribers = make([]chan *Message, 0)
		e.subscribersByTopic[topic] = subscribers
	}

	s := make(chan *Message, capacity)
	subscribers = append(subscribers, s)
	e.subscribersByTopic[topic] = subscribers

	return s
}

func (e *Engine) Unsubscribe(topic string, s <-chan *Message) {
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

func (e *Engine) broadcast(topic string, m *Message) {
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
