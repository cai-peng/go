package dbutils

import (
	"database/sql"
	"strings"
)

type Conn interface {
	Connection() (*sql.DB, error)
}

func NewConn(host, port, user, password, dbname string) Conn {
	return &db{host, port, user, password, dbname}
}

type db struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func (d *db) Connection() (db *sql.DB, err error) {
	if d.dbName == "" {
		d.dbName = "mysql"
	}

	config := []string{d.user, ":",
		d.password, "@tcp(",
		d.host, ":",
		d.port, ")/",
		d.dbName, "?charset=utf8mb4",
		"&timeout=5s&interpolateParams=true",
	}
	db, err = sql.Open("mysql", strings.Join(config, ""))
	if err != nil {
		return nil, err
	}

	//ping这里是必要的因为在Open时即使反馈了异常有时会出现err==nil的情况
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return
}
