package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type PlmnCrawlConf struct {
	Plmn struct {
		Env       string `yaml:"env"`
		URL       string `yaml:"url"`
		LocalFile string `yaml:"localFile"`
		Path      struct {
			Tr string `yaml:"tr"`
			Td string `yaml:"td"`
		} `yaml:"path"`
	} `yaml:"plmn"`
}

func ReadPlmnCrawlConf(path string) *PlmnCrawlConf {
	pc := &PlmnCrawlConf{}
	c, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(c, pc); err != nil {
		log.Fatal(err)
	}

	return pc
}
