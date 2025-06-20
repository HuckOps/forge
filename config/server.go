package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var ServerConfig serverConfig

type serverConfig struct {
	Port     int    `yaml:"port" default:"8080"`
	Host     string `yaml:"host" default:"0.0.0.0"`
	MongoURL string `yaml:"mongo_url" default:"mongodb://localhost:27017"`
}

func InitServerConfig(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	ServerConfig = serverConfig{}
	err = yaml.Unmarshal(data, &ServerConfig)
	if err != nil {
		panic(err)
	}
}
