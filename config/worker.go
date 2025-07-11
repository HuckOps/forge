package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var WorkerConfig workerConfig

type workerConfig struct {
	AMQP      string `yaml:"amqp" default:"amqp://guest@localhost:5672/"`
	MongoURL  string `yaml:"mongo_url" default:"mongodb://localhost:27017"`
	ConfigDir string `yaml:"config_dir" default:"./static"`
}

func InitWorkerConfig(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &WorkerConfig)
	if err != nil {
		panic(err)
	}
}
