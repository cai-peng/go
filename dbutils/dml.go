package dbutils

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	//github.com/outbrain/golib/dbutils
)

type KV map[string]interface{}

type DML struct {
	DB    *sql.DB
	Table string
	*KV
	Where string
}

func (dml *DML) Insert() (sql.Result, error) {
	split := make([]string, 0)
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	for col := range *dml.KV {
		split = append(split, "?")
		columns = append(columns, col)
	}
	//sort map generate a fixed too much type of sql
	sort.Strings(columns)
	for _, v := range columns {
		values = append(values, (*dml.KV)[v])
	}

	s := fmt.Sprint("insert into ", dml.Table, "(", strings.Join(columns, ","), ")values(",
		strings.Join(split, ","), ")")
	return dml.DB.Exec(s, values...)
}

func (dml *DML) Replace() (sql.Result, error) {
	split := make([]string, 0)
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	for col := range *dml.KV {
		split = append(split, "?")
		columns = append(columns, col)
	}

	sort.Strings(columns)
	for _, v := range columns {
		values = append(values, (*dml.KV)[v])
	}

	s := fmt.Sprint("replace into ", dml.Table, "(", strings.Join(columns, ","), ")values(",
		strings.Join(split, ","), ")")
	return dml.DB.Exec(s, values...)
}

func (dml *DML) Update() (sql.Result, error) {
	columns := make([]string, 0)
	values := make([]interface{}, 0)

	keys := make([]string, 0)
	for k := range *dml.KV {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, k := range keys {
		columns = append(columns, fmt.Sprintln(k, "=?"))
		values = append(values, (*dml.KV)[k])
	}

	s := fmt.Sprint("update ", dml.Table, " set ", strings.Join(columns, ","), " where ", dml.Where)
	return dml.DB.Exec(s, values...)
}

func (dml *DML) Delete() (sql.Result, error) {
	s := fmt.Sprint("delete from ", dml.Table, " where ", dml.Where)
	return dml.DB.Exec(s)
}
