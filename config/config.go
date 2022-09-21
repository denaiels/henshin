package config

import (
	"source.golabs.io/go-food/xtools/xconfig"
)

type Config struct {
	AppName     string `yaml:"APP_NAME" env:"APP_NAME"`
	ServerHost  string `yaml:"SERVER_HOST" env:"SERVER_HOST"`
	ServerPort  string `yaml:"SERVER_PORT" env:"SERVER_PORT"`
	StorageHost string `yaml:"STORAGE_HOST" env:"STORAGE_HOST"`
	StoragePort string `yaml:"STORAGE_PORT" env:"STORAGE_PORT"`
}

func NewConfig(filename string) (*Config, error) {
	cfg := &Config{}

	err := xconfig.LoadConfig(filename, cfg)

	if err != nil {
		return nil, err
	}

	return cfg, err
}
