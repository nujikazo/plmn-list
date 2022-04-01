package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/nujikazo/plmn-list/crawl"
	"github.com/nujikazo/plmn-list/crawl/config"

	_ "github.com/mattn/go-sqlite3"
)

const (
	plmn        = "plmn"
	mcc         = "mcc"
	mnc         = "mnc"
	iso         = "iso"
	country     = "country"
	countryCode = "country_code"
	network     = "network"
)

type DB struct {
	*sql.DB
}

// New
func New(conf *config.GeneralConf) (*DB, error) {
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

	stmt := fmt.Sprintf("CREATE TABLE %s (%s text, %s text, %s text, %s text, %s text, %s text); DELETE FROM %s;",
		plmn, mcc, mnc, iso, country, countryCode, network, plmn)

	_, err = db.Exec(stmt)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Insert
func (d *DB) Insert(list []crawl.Plmn) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}

	query, args := d.createBulkInsertQuery(list, 0)

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
func (d *DB) createBulkInsertQuery(list []crawl.Plmn, start int) (query string, args []interface{}) {
	n := len(list)
	values := make([]string, n)
	args = make([]interface{}, n*6)
	pos := 0
	for i := 0; i < n; i++ {
		values[i] = "(?, ?, ?, ?, ?, ?)"
		args[pos] = list[i].MCC
		args[pos+1] = list[i].MNC
		args[pos+2] = list[i].ISO
		args[pos+3] = list[i].Country
		args[pos+4] = list[i].CountryCode
		args[pos+5] = list[i].Network
		pos += 6
	}

	query = fmt.Sprintf(
		"INSERT INTO %s(%s, %s, %s, %s, %s, %s) VALUES %s",
		plmn, mcc, mnc, iso, country, countryCode, network,
		strings.Join(values, ", "),
	)
	return
}
