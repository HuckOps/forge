package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var GatewayConfig gatewayConfig

type gatewayConfig struct {
	Host         string   `yaml:"host" default:"0.0.0.0"`
	Port         int      `yaml:"port" default:"8081"`
	MongoURL     string   `yaml:"mongo_url" default:"mongodb://localhost:27017"`
	AMQP         string   `yaml:"amqp" default:"amqp://guest:guest@localhost:5672/"`
	ConfigDir    string   `yaml:"config_dir" default:"./config"`
	RedisCluster []string `yaml:"redis_cluster" default:"redis://localhost:6379"`
}

func InitGatewayConfig(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	GatewayConfig = gatewayConfig{}
	err = yaml.Unmarshal(data, &GatewayConfig)
	if err != nil {
		panic(err)
	}
}
