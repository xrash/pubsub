package pubsub

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine(t *testing.T) {
	e := NewEngine(1024)
	go e.Start()

	s1 := e.Subscribe("test", 1024)
	s2 := e.Subscribe("test", 1024)
	s3 := e.Subscribe("test", 1024)

	e.Publish("test", "whatever1", 15)
	e.Publish("test", "whatever2", 91)

	testSubscriber := func(s <-chan *Message, name string, data interface{}) {
		m := <-s
		assert.Equal(t, name, m.Name, "Message name should be %v", name)
		assert.Equal(t, data, m.Data, "Message data should be %v", data)
	}

	testSubscriber(s1, "whatever1", 15)
	testSubscriber(s1, "whatever2", 91)

	testSubscriber(s2, "whatever1", 15)
	testSubscriber(s2, "whatever2", 91)

	testSubscriber(s3, "whatever1", 15)
	testSubscriber(s3, "whatever2", 91)
}
