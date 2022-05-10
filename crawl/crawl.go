package main

import (
	"log"
	"os"

	"github.com/nujikazo/plmn-list/config"
	cg "github.com/nujikazo/plmn-list/crawl/config"
	"github.com/nujikazo/plmn-list/crawl/scrape"
	"github.com/nujikazo/plmn-list/database"
)

func main() {
	generalConf := config.New(os.Getenv("GENERAL_CONF"))
	crawlerConf := cg.New(os.Getenv("CRAWLER_CONF"))
	if err := run(generalConf, crawlerConf); err != nil {
		log.Fatal(err)
	}
}

func run(generalConf *config.GeneralConf, crawlerConf *cg.CrawlConf) error {
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

	db.Result = make([]*database.Schema, len(list))
	for i, v := range list {
		db.Result[i] = v
	}

	if err := db.Insert(); err != nil {
		return err
	}

	return nil
}
