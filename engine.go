package pubsub

import (
	"sync"
)

type Engine struct {
	messages chan *wrappedMessage

	topicSubscribersLock *sync.Mutex
	topicSubscribers     map[string][]chan *Message
}

func NewEngine(capacity int) *Engine {
	return &Engine{
		messages: make(chan *wrappedMessage, capacity),

		topicSubscribersLock: &sync.Mutex{},
		topicSubscribers:     make(map[string][]chan *Message),
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
	e.topicSubscribersLock.Lock()
	defer e.topicSubscribersLock.Unlock()

	var list []chan *Message
	var has bool

	list, has = e.topicSubscribers[topic]
	if !has {
		list = make([]chan *Message, 0)
		e.topicSubscribers[topic] = list
	}

	s := make(chan *Message, capacity)
	list = append(list, s)
	e.topicSubscribers[topic] = list

	return s
}

func (e *Engine) Unsubscribe(topic string, s <-chan *Message) {
	e.topicSubscribersLock.Lock()
	defer e.topicSubscribersLock.Unlock()

	list, has := e.topicSubscribers[topic]
	if !has {
		return
	}

	indexToRemove := -1

	for i, sub := range list {
		if sub == s {
			indexToRemove = i
		}
	}

	if indexToRemove < 0 {
		return
	}

	list = append(list[:indexToRemove], list[indexToRemove+1:]...)

	e.topicSubscribers[topic] = list
}

func (e *Engine) broadcast(topic string, m *Message) {
	e.topicSubscribersLock.Lock()
	defer e.topicSubscribersLock.Unlock()

	list, ok := e.topicSubscribers[topic]
	if !ok {
		return
	}

	for _, s := range list {
		s <- m
	}
}
