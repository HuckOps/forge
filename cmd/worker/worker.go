package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/HuckOps/forge/config"
	"github.com/HuckOps/forge/db"
	"github.com/HuckOps/forge/mq"
	consumer2 "github.com/HuckOps/forge/worker/consumer"
	"github.com/HuckOps/forge/worker/cron"
	"os"
)

var (
	consumer   bool
	genConfig  bool
	configPath string
)

func main() {
	flag.StringVar(&configPath, "config", "./config/worker.yaml", "config path")
	flag.BoolVar(&genConfig, "g", false, "generate config file")
	flag.BoolVar(&consumer, "c", false, "consumer mode")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if !genConfig && !consumer {
		flag.Usage()
		os.Exit(1)
	}

	config.InitWorkerConfig(configPath)
	err := mq.InitAMQP(ctx, config.WorkerConfig.AMQP)
	if err != nil {
		panic(err)
	}
	db.InitMongo(ctx, config.WorkerConfig.MongoURL)
	if consumer {
		startConsumer(ctx)
	}
	if genConfig {
		startCron()
	}
	select {}
}

func startConsumer(ctx context.Context) {
	go mq.RabbitMQClient.StartConsumer(ctx, mq.REGISTRYCHANNEL, consumer2.Registry)
	go mq.RabbitMQClient.StartConsumer(ctx, mq.HEARTBEAT, consumer2.Heartbeat)
}

func startCron() {
	cron.StartGenGatewayConfig()
}
