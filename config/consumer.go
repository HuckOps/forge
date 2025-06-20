package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var ConsumerConfig consumerConfig

type consumerConfig struct {
	AMQP     string `yaml:"amqp" default:"amqp://guest@localhost:5672/"`
	MongoURL string `yaml:"mongo_url" default:"mongodb://localhost:27017"`
}

func InitConsumerConfig(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &ConsumerConfig)
	if err != nil {
		panic(err)
	}
}
