package mq

import (
	"context"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	lock    sync.Mutex
}

var RabbitMQClient *rabbitMQClient

const (
	REGISTRYCHANNEL   = "registry"
	UNREGISTRYCHANNEL = "unregistry"
	HEARTBEAT         = "heartbeat"
)

type HeartBeatMessage struct {
	UUID string `json:"uuid"`
}

type RegisteyMessage struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	UUID     string `json:"uuid"`
}

func InitAMQP(ctx context.Context, amqpURL string) error {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		return err
	}

	queues := []string{REGISTRYCHANNEL, UNREGISTRYCHANNEL, HEARTBEAT}
	for _, q := range queues {
		_, err := channel.QueueDeclare(
			q,
			false, // durable
			false, // autoDelete
			false, // exclusive
			false, // noWait
			nil,
		)
		if err != nil {
			return err
		}
	}

	RabbitMQClient = &rabbitMQClient{
		conn:    conn,
		channel: channel,
	}

	log.Println("RabbitMQ connection established.")
	return nil
}

func (c *rabbitMQClient) StartConsumer(ctx context.Context, queue string, handlerFunc ...func(amqp.Delivery)) {

	for {
		select {
		case <-ctx.Done():
			log.Printf("Stopping consumer for queue: %s", queue)
			return
		default:
		}

		msgs, err := c.channel.Consume(
			queue,
			"",
			false, // auto-ack disabled
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Println("Consume failed:", err)
			time.Sleep(time.Second)
			continue
		}

		for d := range msgs {
			go func(delivery amqp.Delivery) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered in handler: %v", r)
					}
				}()
				for _, h := range handlerFunc {
					h(delivery)
				}
				delivery.Ack(false)
				log.Printf("Handled message from %s: %s", queue, delivery.Body)
			}(d)
		}
	}

}

func (c *rabbitMQClient) PublishMessage(queueName string, body []byte, persistent bool) error {
	if c.channel == nil {
		return amqp.ErrClosed
	}

	// 设定 deliveryMode
	deliveryMode := amqp.Transient
	if persistent {
		deliveryMode = amqp.Persistent
	}

	err := c.channel.Publish(
		"",        // exchange: 空表示默认
		queueName, // routing key（即 queue 名）
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         body,
			DeliveryMode: deliveryMode, // 1=non-persistent, 2=persistent
		},
	)
	if err != nil {
		log.Printf("Failed to publish message to queue %s: %v", queueName, err)
	}
	return err
}

func (c *rabbitMQClient) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	log.Println("RabbitMQ connection closed.")
}
