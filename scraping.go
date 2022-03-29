package main

import (
	"log"

	"github.com/nujikazo/plmn-list/crawl"
	"github.com/nujikazo/plmn-list/crawl/config"
)

func main() {
	conf := config.ReadGeneralConf("")
	if err := crawl.Run(conf); err != nil {
		log.Fatal(err)
	}
}
