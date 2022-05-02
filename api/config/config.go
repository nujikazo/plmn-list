package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type APIConf struct {
	ServerAddr  string `yaml:"serverAddr"`
	ServerPort  int    `yaml:"serverPort"`
	GatewayAddr string `yaml:"gatewayAddr"`
	GatewayPort int    `yaml:"gatewayPort"`
}

func New(path string) *APIConf {
	ac := &APIConf{}
	c, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(c, ac); err != nil {
		log.Fatal(err)
	}

	return ac
}
