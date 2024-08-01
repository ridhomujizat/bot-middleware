package messaging

type MessagingPublisher interface {
	Publish(queueName string, payload interface{}) error
}

type MessagingSubscriber interface {
	Subscribe(queueName string, handleFunc func([]byte)) error
}

type MessagingGeneral interface {
	GetPublisher() MessagingPublisher
	GetSubscriber() MessagingSubscriber
}

type messagingGeneral struct {
	publisher  MessagingPublisher
	subscriber MessagingSubscriber
}

func NewMessagingGeneral(publisher MessagingPublisher, subscriber MessagingSubscriber) MessagingGeneral {
	return &messagingGeneral{publisher: publisher, subscriber: subscriber}
}

func (m *messagingGeneral) GetPublisher() MessagingPublisher {
	return m.publisher
}

func (m *messagingGeneral) GetSubscriber() MessagingSubscriber {
	return m.subscriber
}
