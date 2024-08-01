package rabbit

import (
	"bot-middleware/config"

	"github.com/pterm/pterm"
	"github.com/streadway/amqp"
)

type RabbitMQSubscriber struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQSubscriber(cfg config.RabbitMQConfig) (*RabbitMQSubscriber, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQSubscriber{
		connection: conn,
		channel:    ch,
	}, nil
}

func (r *RabbitMQSubscriber) Subscribe(queueName string, handleFunc func([]byte)) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		queueName,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handleFunc(msg.Body)
		}
	}()

	pterm.Info.Printfln(" [*] Waiting for messages in %s", queueName)
	select {}
}

func (r *RabbitMQSubscriber) Close() {
	r.channel.Close()
	r.connection.Close()
}
