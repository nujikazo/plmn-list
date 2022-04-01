package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/nujikazo/plmn-list/crawl/config"

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
func New(conf *config.GeneralConf) (*Database, error) {
	target := fmt.Sprintf("%s", conf.DatabaseName)

	_, err := os.Stat(target)

	if !os.IsNotExist(err) {
		if err := os.Remove(target); err != nil {
			return nil, err
		}
	}

	db, err := sql.Open(conf.DatabaseType, target)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf("CREATE TABLE %s (%s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT); DELETE FROM %s;",
		table, mcc, mnc, iso, country, countryCode, network, table)

	_, err = db.Exec(stmt)
	if err != nil {
		return nil, err
	}

	var schema []Schema
	return &Database{
		db,
		schema,
	}, nil
}

// Insert
func (db *Database) Insert() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	start := 0
	query, args := db.createBulkInsertQuery(start)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(args...); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

// createBulkInsertQuery
func (db *Database) createBulkInsertQuery(start int) (string, []interface{}) {
	n := len(db.Schemas)
	values := make([]string, n)
	args := make([]interface{}, n*6)
	pos := 0
	for i := 0; i < n; i++ {
		values[i] = "(?, ?, ?, ?, ?, ?)"
		args[pos] = db.Schemas[i].MCC
		args[pos+1] = db.Schemas[i].MNC
		args[pos+2] = db.Schemas[i].ISO
		args[pos+3] = db.Schemas[i].Country
		args[pos+4] = db.Schemas[i].CountryCode
		args[pos+5] = db.Schemas[i].Network
		pos += 6
	}

	query := fmt.Sprintf(
		"INSERT INTO %s(%s, %s, %s, %s, %s, %s) VALUES %s",
		table, mcc, mnc, iso, country, countryCode, network,
		strings.Join(values, ", "),
	)

	return query, args
}
