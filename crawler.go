package main

import (
	"log"
	"os"

	"github.com/nujikazo/plmn-list/crawl"
	"github.com/nujikazo/plmn-list/crawl/config"
)

func main() {
	generalConf := config.ReadGeneralConf(os.Getenv("GENERAL_CONF"))
	crawlerConf := config.ReadPlmnCrawlConf(os.Getenv("CRAWLER_CONF"))
	if err := crawl.Run(generalConf, crawlerConf); err != nil {
		log.Fatal(err)
	}
}
