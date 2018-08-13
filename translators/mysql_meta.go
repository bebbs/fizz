package translators

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gobuffalo/fizz"
	"github.com/pkg/errors"
)

type mysqlTableInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default interface{}
	Extra   string
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
	db, err := sql.Open("mysql", p.URL)
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("select TABLE_NAME as name from information_schema.TABLES where TABLE_SCHEMA = '%s'", p.Name))
	if err != nil {
		return err
	}
	defer res.Close()
	for res.Next() {
		table := &fizz.Table{
			Columns: []fizz.Column{},
			Indexes: []fizz.Index{},
		}
		err = res.Scan(&table.Name)
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

func (p *mysqlSchema) buildTableData(table *fizz.Table, db *sql.DB) error {
	prag := fmt.Sprintf("SELECT `COLUMN_NAME` AS `Field`, `COLUMN_TYPE` AS `Type`, `IS_NULLABLE` AS `Null`, `COLUMN_KEY` AS `Key`, `COLUMN_DEFAULT` AS `Default`, `EXTRA` AS `Extra` FROM information_schema.columns WHERE table_name='%s';", table.Name)

	res, err := db.Query(prag)
	if err != nil {
		return nil
	}

	for res.Next() {
		ti := mysqlTableInfo{}
		err = res.Scan(&ti.Field, &ti.Type, &ti.Null, &ti.Key, &ti.Default, &ti.Extra)
		if err != nil {
			return err
		}
		table.Columns = append(table.Columns, ti.ToColumn())
	}

	p.schema[table.Name] = table
	return nil
}
