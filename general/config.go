package general

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// GeneralConf
type GeneralConf struct {
	DatabaseHost     string `yaml:"databaseHost"`
	DatabaseUser     string `yaml:"databaseUser"`
	DatabasePassword string `yaml:"databasePassword"`
	DatabasePort     int    `yaml:"databasePort"`
	DatabaseName     string `yaml:"databaseName"`
	DatabassCharset  string `yaml:"databassCharset"`
	DatabaseType     string `yaml:"databaseType"`
}

// ReadGeneralConf
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