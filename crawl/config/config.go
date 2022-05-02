package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type CrawlConf struct {
	Plmn []struct {
		Env       string `yaml:"env"`
		URL       string `yaml:"url"`
		LocalFile string `yaml:"localFile"`
		Path      struct {
			Tr string `yaml:"tr"`
			Td string `yaml:"td"`
		} `yaml:"path"`
	} `yaml:"plmn"`
}

func New(path string) *CrawlConf {
	pc := &CrawlConf{}
	c, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(c, pc); err != nil {
		log.Fatal(err)
	}

	return pc
}
