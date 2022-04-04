package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type APIConf struct {
	ServerAddr string `yaml:"serverAddr"`
	ServerPort int    `yaml:"serverPort"`
}

func ReadAPIConf(path string) *APIConf {
	gc := &APIConf{}
	c, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(c, gc); err != nil {
		log.Fatal(err)
	}

	return gc
}
