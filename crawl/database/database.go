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

	stmt := `
	CREATE TABLE plmn (mcc text, mnc text, iso text, country text, country_code text, network text);
	DELETE FROM plmn;
	`
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
		"INSERT INTO plmn(mcc, mnc, iso, country, country_code, network) VALUES %s",
		strings.Join(values, ", "),
	)
	return
}
