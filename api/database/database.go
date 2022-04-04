package database

import (
	"database/sql"
	"fmt"

	"github.com/nujikazo/plmn-list/general"

	_ "github.com/mattn/go-sqlite3"
)

const (
	table       = "plmn"
	mcc         = "mcc"
	mnc         = "mnc"
	iso         = "iso"
	country     = "country"
	countryCode = "country_code"
	network     = "network"
)

type Database struct {
	*sql.DB
	Schemas []Schema
}

type Schema struct {
	MCC         string
	MNC         string
	ISO         string
	Country     string
	CountryCode string
	Network     string
}

// New
func New(conf *general.GeneralConf) (*Database, error) {
	target := fmt.Sprintf("%s", conf.DatabaseName)

	db, err := sql.Open(conf.DatabaseType, target)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	var schema []Schema
	return &Database{
		db,
		schema,
	}, nil
}

// Insert
func (db *Database) GetPlmnList() ([]Schema, error) {
	stmt := fmt.Sprintf("SELECT * FROM %s;", general.Table)
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Schema

	for rows.Next() {
		var plmn Schema
		if err := rows.Scan(
			&plmn.MCC,
			&plmn.MNC,
			&plmn.ISO,
			&plmn.Country,
			&plmn.CountryCode,
			&plmn.Network,
		); err != nil {
			return nil, err
		}

		list = append(list, plmn)
	}

	return list, nil
}
