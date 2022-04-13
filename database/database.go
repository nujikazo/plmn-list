package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/nujikazo/plmn-list/general"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
	Schemas []Schema
	Name    string
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
		target,
	}, nil
}

// InitializeDB
func (db *Database) InitializeDB() error {
	stmt := fmt.Sprintf("CREATE TABLE %s (%s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT, %s TEXT); DELETE FROM %s;",
		general.Table,
		general.Mcc,
		general.Mnc,
		general.Iso,
		general.Country,
		general.CountryCode,
		general.Network,
		general.Table,
	)

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
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
		general.Table,
		general.Mcc,
		general.Mnc,
		general.Iso,
		general.Country,
		general.CountryCode,
		general.Network,
		strings.Join(values, ", "),
	)

	return query, args
}

// GetPlmnList
func (db *Database) GetPlmnList(query map[string]string) ([]Schema, error) {
	var args []interface{}
	var values string
	var list []Schema
	stmt := fmt.Sprintf("SELECT * FROM %s", general.Table)

	if db.checkQuery(query) {
		values, args = db.buildGetQuery(query)
		stmt = fmt.Sprintf("%s %s", stmt, values)
	}

	rows, err := db.Query(stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

// buildGetQuery
func (db *Database) buildGetQuery(query map[string]string) (string, []interface{}) {
	var result []string
	var args []interface{}

	for k, v := range query {
		if v != "" {
			values := fmt.Sprintf("%s = ?", k)
			result = append(result, values)
			args = append(args, v)
		}
	}

	return fmt.Sprintf("WHERE %s;", strings.Join(result, " AND ")), args
}

// checkQuery
func (db *Database) checkQuery(query map[string]string) bool {
	var check = false
	for _, v := range query {
		if v != "" {
			check = true
		}
	}

	return check
}
