package messaging

type Publisher interface {
	Publish(queueName string, payload interface{}) error
}

type Subscriber interface {
	Subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool, handleFunc func([]byte) error) error
}

type MessagingGeneral interface {
	Publish(queueName string, payload interface{}) error
	Subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool, handleFunc func([]byte) error) error
}

type messagingGeneral struct {
	publisher  Publisher
	subscriber Subscriber
}

func NewMessagingGeneral(publisher Publisher, subscriber Subscriber) MessagingGeneral {
	return &messagingGeneral{publisher: publisher, subscriber: subscriber}
}

func (m *messagingGeneral) Publish(queueName string, payload interface{}) error {
	return m.publisher.Publish(queueName, payload)
}

func (m *messagingGeneral) Subscribe(exchange, routingKey, queueName string, allowNonJsonMessages bool, handleFunc func([]byte) error) error {
	return m.subscriber.Subscribe(exchange, routingKey, queueName, allowNonJsonMessages, handleFunc)
}
