package logic

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DBConnect(uri string) (conn *sql.DB, err error) {
	conn, err = sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	return conn, err
}
