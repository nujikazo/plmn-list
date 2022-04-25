package main

import (
	"log"
	"os"

	"github.com/nujikazo/plmn-list/crawl/config"
	"github.com/nujikazo/plmn-list/crawl/scrape"
	"github.com/nujikazo/plmn-list/database"
	"github.com/nujikazo/plmn-list/general"
)

func main() {
	generalConf := general.ReadGeneralConf(os.Getenv("GENERAL_CONF"))
	crawlerConf := config.ReadPlmnCrawlConf(os.Getenv("CRAWLER_CONF"))
	if err := run(generalConf, crawlerConf); err != nil {
		log.Fatal(err)
	}
}

func run(generalConf *general.GeneralConf, crawlerConf *config.PlmnCrawlConf) error {
	_, err := os.Stat(generalConf.DatabaseName)
	if !os.IsNotExist(err) {
		if err := os.Remove(generalConf.DatabaseName); err != nil {
			return err
		}
	}

	db, err := database.New(generalConf)
	if err != nil {
		return err
	}

	if err := db.InitializeDB(); err != nil {
		return err
	}

	list, err := scrape.Run(generalConf, crawlerConf)
	if err != nil {
		return err
	}

	db.Result = make([]database.Schema, len(list))
	for i, v := range list {
		db.Result[i].MCC = v.MCC
		db.Result[i].MNC = v.MNC
		db.Result[i].ISO = v.ISO
		db.Result[i].Country = v.Country
		db.Result[i].CountryCode = v.CountryCode
		db.Result[i].Network = v.Network
	}

	if err := db.Insert(); err != nil {
		return err
	}

	return nil
}
