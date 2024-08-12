package rabbit

import (
	"bot-middleware/config"
	"bot-middleware/internal/pkg/repository/redis"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/pterm/pterm"
	redisClient "github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
)

type RabbitMQSubscriber struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	redis      *redis.RedisClient
}

func NewRabbitMQSubscriber(cfg config.RabbitMQConfig, redis *redis.RedisClient) (*RabbitMQSubscriber, error) {
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
		redis:      redis,
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
				headerId, ok := msg.Headers["id"].(string)
				if !ok {
					log.Println("Message missing 'id' header, skipping...")
					msg.Ack(false)
					continue
				}

				retryKey := fmt.Sprintf("retry:%s", headerId)
				retryCount, err := r.getRetryCount(retryKey)
				if err != nil {
					pterm.Error.Printfln("Failed to get retry count for message %s: %v", headerId, err)
					msg.Nack(false, true)
					continue
				}

				if retryCount >= 3 {
					pterm.Error.Printfln("Message %s exceeded retry limit, discarding...", headerId)
					msg.Ack(false)
					r.deleteRetryCount(retryKey)
					continue
				}

				err = handleFunc(msg.Body)
				if err != nil {
					pterm.Warning.Printfln("Failed to process message: %v. Requeuing...", err)
					r.incrementRetryCount(retryKey)
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

func (r *RabbitMQSubscriber) getRetryCount(key string) (int, error) {
	val, err := r.redis.Get(key)
	if err == redisClient.Nil {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}

	retryCount, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return retryCount, nil
}

func (r *RabbitMQSubscriber) incrementRetryCount(key string) error {
	val, err := r.redis.Get(key)
	if err != nil && err != redisClient.Nil {
		return err
	}

	retryCount := 0

	if val != "" {
		retryCount, err = strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("failed to parse retry count for key %s: %w", key, err)
		}
	}

	retryCount++

	err = r.redis.Set(key, retryCount, 5*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to set retry count for key %s: %w", key, err)
	}

	return nil
}

func (r *RabbitMQSubscriber) deleteRetryCount(key string) {
	r.redis.Del(key)
}

func (r *RabbitMQSubscriber) Close() {
	r.channel.Close()
	r.connection.Close()
}
