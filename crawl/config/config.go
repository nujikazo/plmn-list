package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type GeneralConf struct {
	DatabaseHost     string `yaml:"databaseHost"`
	DatabaseUser     string `yaml:"databaseUser"`
	DatabasePassword string `yaml:"databasePassword"`
	DatabasePort     int    `yaml:"databasePort"`
	DatabaseName     string `yaml:"databaseName"`
	DatabassCharset  string `yaml:"databassCharset"`
	DatabaseType     string `yaml:"databaseType"`
}

func ReadGeneralConf(path string) *GeneralConf {
	gc := &GeneralConf{}
	c, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(c, gc); err != nil {
		log.Fatal(err)
	}

	return gc
}

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
