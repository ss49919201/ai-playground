package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func NewClient(
	user,
	password,
	host,
	database string,
) (*sql.DB, error) {
	config := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   host,
		DBName: database,
	}

	return sql.Open("mysql", config.FormatDSN())
}
