package rabbit

import (
	"bot-middleware/config"
	"bot-middleware/internal/pkg/repository/redis"
	"bot-middleware/internal/pkg/util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"
	redisClient "github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
)

type RabbitMQSubscriber struct {
	config     config.RabbitMQConfig
	connection *amqp.Connection
	channel    *amqp.Channel
	redis      *redis.RedisClient
	publisher  *RabbitMQPublisher
}

func NewRabbitMQSubscriber(cfg config.RabbitMQConfig, publisher *RabbitMQPublisher, redis *redis.RedisClient) (*RabbitMQSubscriber, error) {
	subscriber := &RabbitMQSubscriber{
		config:    cfg,
		publisher: publisher,
		redis:     redis,
	}

	err := subscriber.connect()
	if err != nil {
		return nil, err
	}

	go subscriber.handleConnectionClosure()

	return subscriber, nil
}

func (r *RabbitMQSubscriber) connect() error {
	var err error
	for {
		r.connection, err = amqp.Dial(r.config.URL)
		if err != nil {
			pterm.Warning.Printfln("Failed to connect to RabbitMQ, retrying in 5 seconds: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		r.channel, err = r.connection.Channel()
		if err != nil {
			pterm.Warning.Printfln("Failed to open a channel, retrying in 5 seconds: %v", err)
			r.connection.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		pterm.Info.Println("Connected to RabbitMQ")
		return nil
	}
}

func (r *RabbitMQSubscriber) handleConnectionClosure() {
	closeCh := r.connection.NotifyClose(make(chan *amqp.Error))
	for err := range closeCh {
		if err != nil {
			pterm.Error.Printfln("Connection closed: %v. Reconnecting...", err)
			r.reconnect()
		}
	}
}

func (r *RabbitMQSubscriber) reconnect() {
	r.Close()
	time.Sleep(5 * time.Second)
	err := r.connect()
	if err != nil {
		pterm.Fatal.Printfln("Failed to reconnect to RabbitMQ: %v", err)
	}
}

func (r *RabbitMQSubscriber) Subscribe(exchange, routingKey string, queueName string, allowNonJsonMessages bool, handleFunc func([]byte) error) error {
	for {
		if r.connection == nil || r.channel == nil {
			if err := r.connect(); err != nil {
				return err
			}
		}

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
			pterm.Error.Printfln("ExchangeDeclare failed: %v. Reconnecting...", err)
			r.reconnect()
			continue
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
			pterm.Error.Printfln("QueueDeclare failed: %v. Reconnecting...", err)
			r.reconnect()
			continue
		}

		err = r.channel.QueueBind(
			q.Name,
			routingKey,
			exchange,
			false,
			nil,
		)
		if err != nil {
			pterm.Error.Printfln("QueueBind failed: %v. Reconnecting...", err)
			r.reconnect()
			continue
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
			pterm.Error.Printfln("Consume failed: %v. Reconnecting...", err)
			r.reconnect()
			continue
		}

		go r.handleMessages(msgs, allowNonJsonMessages, queueName, handleFunc)

		pterm.Info.Printfln(" [*] Waiting for messages in %s", queueName)
		select {}
	}
}

func (r *RabbitMQSubscriber) handleMessages(msgs <-chan amqp.Delivery, allowNonJsonMessages bool, queueName string, handleFunc func([]byte) error) {
	for msg := range msgs {
		if allowNonJsonMessages || json.Valid(msg.Body) {
			headerId, ok := msg.Headers["id"].(string)
			if !ok {
				pterm.Warning.Println("Message missing 'id' header, skipping...")
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
				pterm.Error.Printfln("Message %s exceeded retry limit, moving to fail...", headerId)

				failQueueName := queueName
				if !strings.HasPrefix(queueName, "fail:") {
					failQueueName = "fail:" + queueName
				}
				var payload interface{}
				err := json.Unmarshal(msg.Body, &payload)
				if err != nil {
					pterm.Error.Printfln("Unmarshal failed: %v", err)
					msg.Nack(false, true)
					continue
				}

				err = r.publisher.Publish(failQueueName, payload)
				if err != nil {
					pterm.Error.Printfln("Publish failed: %v", err)
					msg.Nack(false, true)
					continue
				}

				msg.Ack(false)
				r.deleteRetryCount(retryKey)
				continue
			}

			err = handleFunc(msg.Body)
			if err != nil {
				pterm.Warning.Printfln("Failed to process message: %v. Requeuing...", err)
				r.incrementRetryCount(retryKey)
				util.Sleep(2000)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		} else {
			pterm.Warning.Printfln("Received non-JSON message: %s", string(msg.Body))
			msg.Nack(false, true)
		}
	}
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
	if r.channel != nil {
		r.channel.Close()
	}
	if r.connection != nil {
		r.connection.Close()
	}
}
