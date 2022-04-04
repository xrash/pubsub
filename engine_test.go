package pubsub

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine(t *testing.T) {
	e := NewEngine[int](1024)
	go e.Start()

	s1 := e.Subscribe("test", 1024)
	s2 := e.Subscribe("test", 1024)
	s3 := e.Subscribe("test", 1024)

	e.Publish("test", 15)
	e.Publish("test", 91)

	testSubscriber := func(s <-chan int, expectedMsg int) {
		m := <-s
		assert.Equal(t, expectedMsg, m, "Message should be %v", expectedMsg)
	}

	testSubscriber(s1, 15)
	testSubscriber(s1, 91)

	testSubscriber(s2, 15)
	testSubscriber(s2, 91)

	testSubscriber(s3, 15)
	testSubscriber(s3, 91)
}
