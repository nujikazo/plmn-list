package main

import (
	"log"
	"os"

	"github.com/nujikazo/plmn-list/crawl"
	"github.com/nujikazo/plmn-list/crawl/config"
	"github.com/nujikazo/plmn-list/crawl/database"
)

func main() {
	generalConf := config.ReadGeneralConf(os.Getenv("GENERAL_CONF"))
	crawlerConf := config.ReadPlmnCrawlConf(os.Getenv("CRAWLER_CONF"))
	if err := run(generalConf, crawlerConf); err != nil {
		log.Fatal(err)
	}
}

func run(generalConf *config.GeneralConf, crawlerConf *config.PlmnCrawlConf) error {
	db, err := database.New(generalConf)
	if err != nil {
		return err
	}

	list, err := crawl.Run(generalConf, crawlerConf)
	if err != nil {
		return err
	}

	if err := db.Insert(list); err != nil {
		return err
	}

	return nil
}
