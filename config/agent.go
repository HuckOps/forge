package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

var AgentConfig agentConfig

type agentConfig struct {
	RegistryCenter string `yaml:"registry_center"`
}

func InitAgentConfig(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &AgentConfig)
	if err != nil {
		panic(err)
	}
}
