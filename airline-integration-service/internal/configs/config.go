package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Elasticsearch struct {
		Addresses []string `yaml:"addresses"`
	} `yaml:"elasticsearch"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Log struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
	} `yaml:"log"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{
		Elasticsearch: struct {
			Addresses []string `yaml:"addresses"`
		}{
			Addresses: []string{"http://localhost:9200"},
		},
		Server: struct {
			Port string `yaml:"port"`
		}{
			Port: "8080",
		},
		Log: struct {
			Path  string `yaml:"path"`
			Level string `yaml:"level"`
		}{
			Path:  "./logs/app.log",
			Level: "debug",
		},
	}

	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %v", err)
		}
	}

	if addr := os.Getenv("ES_ADDRESSES"); addr != "" {
		cfg.Elasticsearch.Addresses = []string{addr}
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}
	if logPath := os.Getenv("LOG_PATH"); logPath != "" {
		cfg.Log.Path = logPath
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.Log.Level = logLevel
	}

	return cfg, nil
}
