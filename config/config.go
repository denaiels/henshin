package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AppName     string `yaml:"APP_NAME"`
	Host        string `yaml:"HOST"`
	ServerPort  string `yaml:"SERVER_PORT"`
	StoragePort string `yaml:"STORAGE_PORT"`
}

func NewConfig(path string) Config {
	var cfg Config

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.NewConfig err    #%v ", err)
	}
	print(yamlFile)
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return cfg
}