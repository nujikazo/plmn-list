package main

import (
	"log"
	"os"

	"github.com/nujikazo/plmn-list/crawl/config"
	"github.com/nujikazo/plmn-list/crawl/database"
	"github.com/nujikazo/plmn-list/crawl/scrape"
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

	list, err := scrape.Run(generalConf, crawlerConf)
	if err != nil {
		return err
	}

	db.Schemas = make([]database.Schema, len(list))
	for i, v := range list {
		db.Schemas[i].MCC = v.MCC
		db.Schemas[i].MNC = v.MNC
		db.Schemas[i].ISO = v.ISO
		db.Schemas[i].Country = v.Country
		db.Schemas[i].CountryCode = v.CountryCode
		db.Schemas[i].Network = v.Network
	}

	if err := db.Insert(); err != nil {
		return err
	}

	return nil
}
