package rabbit

import (
	"bot-middleware/config"
	"bot-middleware/internal/pkg/util"
	"encoding/json"
	"log"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/streadway/amqp"
)

type RabbitMQPublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQPublisher(cfg config.RabbitMQConfig) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQPublisher{
		connection: conn,
		channel:    ch,
	}, nil
}

func (r *RabbitMQPublisher) publish(queueName string, message []byte, headers map[string]interface{}) error {
	err := r.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
			Headers:     headers,
		})
	if err != nil {
		util.HandleAppError(err, "publish", "Publish", false)
		return err
	}

	return nil
}

func (r *RabbitMQPublisher) assertQueue(queueName string) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		util.HandleAppError(err, "assertQueue", "QueueDeclare", false)
	}
	return err
}

func (r *RabbitMQPublisher) Publish(queueName string, data interface{}) error {
	if err := r.assertQueue(queueName); err != nil {
		return err
	}

	message, err := json.Marshal(data)
	if err != nil {
		util.HandleAppError(err, "Publish", "Marshal", false)
		return err
	}

	id, err := gonanoid.New()
	if err != nil {
		log.Fatalf("Failed to generate NanoID: %v", err)
	}

	headers := map[string]interface{}{
		"id": id,
	}

	if err := r.publish(queueName, message, headers); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQPublisher) Close() {
	if err := r.channel.Close(); err != nil {
		util.HandleAppError(err, "Close", "channel.Close", false)
	}
	if err := r.connection.Close(); err != nil {
		util.HandleAppError(err, "Close", "connection.Close", false)
	}
}
