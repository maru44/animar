package helper

import (
	"database/sql"
	"fmt"
)

var MysqlUser = GetenvOrDefault("MYSQL_USER", "go")
var MysqlPassword = GetenvOrDefault("MYSQL_PASSWORD", "Go1234_test")
var MysqlDataBase = GetenvOrDefault("MYSQL_DB", "go_test")

func AccessDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", MysqlUser, MysqlPassword, MysqlDataBase))
	if err != nil {
		panic(err.Error())
	}

	return db
}

func NewNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
