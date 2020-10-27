package dbutils

import (
	"database/sql"
	"strings"
)

func Query(db *sql.DB, sqls string) ([]map[string]string, error) {
	var err error
	rows, err := db.Query(sqls)
	records := make([]map[string]string, 0, 500)
	if err != nil {
		return records, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	val_point := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		val_point[i] = &values[i]
	}

	for rows.Next() {
		rows.Scan(val_point...)
		record := make(map[string]string)
		for i, v := range values {
			if v == nil {
				record[strings.ToLower(columns[i])] = "NULL"
			} else {
				record[strings.ToLower(columns[i])] = string(v.([]byte))
			}
		}
		records = append(records, record)
	}
	return records, err
}

func QueryColumn(db *sql.DB, sqls string) ([][]string, error) {
	var err error
	records := make([][]string, 0, 500)
	rows, err := db.Query(sqls)
	if err != nil {
		return records, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	val_point := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		val_point[i] = &values[i]
	}

	for rows.Next() {
		rows.Scan(val_point...)
		record := make([]string, 0, len(columns))
		for _, v := range values {
			if v == nil {
				record = append(record, "NULL")
			} else {
				record = append(record, string(v.([]byte)))
			}
		}
		records = append(records, record)
	}
	return records, err
}

func QueryOne(db *sql.DB, sqls string) (string, error) {
	rows, err := db.Query(sqls)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var record string
	for rows.Next() {
		rows.Scan(&record)
	}
	return record, nil
}
