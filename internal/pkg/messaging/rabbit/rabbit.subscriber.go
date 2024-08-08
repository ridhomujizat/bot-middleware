package rabbit

import (
	"bot-middleware/config"
	"encoding/json"
	"log"

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
		conn.Close()
		return nil, err
	}

	return &RabbitMQSubscriber{
		connection: conn,
		channel:    ch,
	}, nil
}

func (r *RabbitMQSubscriber) Subscribe(exchange, routingKey string, queueName string, allowNonJsonMessages bool, handleFunc func([]byte) error) error {
	err := r.channel.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	q, err := r.channel.QueueDeclare(
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

	err = r.channel.QueueBind(
		q.Name,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if allowNonJsonMessages || json.Valid(msg.Body) {
				err := handleFunc(msg.Body)
				if err != nil {
					log.Printf("Failed to process message: %v. Requeuing...\n", err)
					msg.Nack(false, true)
				} else {
					msg.Ack(false)
				}
			} else {
				log.Printf("Received non-JSON message: %s", string(msg.Body))
				msg.Nack(false, true)
			}
		}
	}()

	pterm.Info.Printfln(" [*] Waiting for messages in %s", queueName)
	select {}
}

func (r *RabbitMQSubscriber) Close() {
	r.channel.Close()
	r.connection.Close()
}
