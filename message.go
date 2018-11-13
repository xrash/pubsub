package pubsub

type wrappedMessage struct {
	topic   string
	message *Message
}

type Message struct {
	Name string
	Data interface{}
}
