package store

type SubscriptionHandler func(event Event) error
type EventStore interface {
	Publish(topic string, data []byte) error
	Subscribe(topic string, handler SubscriptionHandler) error
}

type Event struct {
	topic string
	data  []byte
}

func (e Event) Data() []byte {
	return e.data
}

func (e Event) Topic() string {
	return e.topic
}
