package db

import (
	"database/sql"
	"fmt"

	"github.com/nujikazo/plmn-list/crawl/config"
	"github.com/nujikazo/plmn-list/database/models"

	_ "github.com/go-sql-driver/mysql"
)

func New(conf *config.GeneralConf) (*models.Queries, error) {
	target := fmt.Sprintf("%s:%s@%s/%s?charset=%s&parseTime=%t",
		conf.DatabaseUser,
		conf.DatabasePassword,
		conf.DatabaseHost,
		conf.DatabaseName,
		conf.DatabassCharset,
		true)
	db, err := sql.Open(conf.DatabaseType, target)
	if err != nil {
		return nil, err
	}

	return models.New(db), nil
}
