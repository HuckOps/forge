package main

import (
	"context"
	"github.com/HuckOps/forge/config"
	"github.com/HuckOps/forge/consumer"
	"github.com/HuckOps/forge/db"
	"github.com/HuckOps/forge/mq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.InitConsumerConfig("./config/consumer.yaml")
	err := mq.InitAMQP(ctx, config.ConsumerConfig.AMQP)
	if err != nil {
		panic(err)
	}
	db.InitMongo(ctx, config.ConsumerConfig.MongoURL)
	go mq.RabbitMQClient.StartConsumer(ctx, mq.REGISTRYCHANNEL, consumer.Registry)
	go mq.RabbitMQClient.StartConsumer(ctx, mq.HEARTBEAT, consumer.Heartbeat)
	select {}
}
