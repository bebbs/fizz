package translators

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gobuffalo/fizz"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type mysqlTableInfo struct {
	Field   string      `db:"Field"`
	Type    string      `db:"Type"`
	Null    string      `db:"Null"`
	Key     string      `db:"Key"`
	Default interface{} `db:"Default"`
	Extra   string      `db:"Extra"`
}

func (ti mysqlTableInfo) ToColumn() fizz.Column {
	c := fizz.Column{
		Name:    ti.Field,
		ColType: ti.Type,
		Primary: ti.Key == "PRI",
		Options: map[string]interface{}{},
	}
	if strings.ToLower(ti.Null) == "yes" {
		c.Options["null"] = true
	}
	if ti.Default != nil {
		d := fmt.Sprintf("%s", ti.Default)
		c.Options["default"] = d
	}
	return c
}

type mysqlSchema struct {
	Schema
}

func (p *mysqlSchema) Version() (string, error) {
	var version string
	var err error

	db, err := sql.Open("mysql", p.URL)
	if err != nil {
		return version, err
	}
	defer db.Close()

	res, err := db.Query("select VERSION()")
	if err != nil {
		return version, err
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&version)
		return version, err
	}
	return "", errors.New("could not locate MySQL version")
}

func (p *mysqlSchema) Build() error {
	var err error
	db, err := sqlx.Open("mysql", p.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Queryx(fmt.Sprintf("select TABLE_NAME as name from information_schema.TABLES where TABLE_SCHEMA = '%s'", p.Name))
	if err != nil {
		return err
	}
	for res.Next() {
		table := &fizz.Table{
			Columns: []fizz.Column{},
			Indexes: []fizz.Index{},
		}
		err = res.StructScan(table)
		if err != nil {
			return err
		}
		err = p.buildTableData(table, db)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *mysqlSchema) buildTableData(table *fizz.Table, db *sqlx.DB) error {
	prag := fmt.Sprintf("describe %s", table.Name)

	res, err := db.Queryx(prag)
	if err != nil {
		return nil
	}

	for res.Next() {
		ti := mysqlTableInfo{}
		err = res.StructScan(&ti)
		if err != nil {
			return err
		}
		table.Columns = append(table.Columns, ti.ToColumn())
	}

	p.schema[table.Name] = table
	return nil
}
