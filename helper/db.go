package helper

import (
	"database/sql"
	"fmt"
)

type NullInt struct {
	Int   int
	Valid bool
}

var (
	MysqlUser     = GetenvOrDefault("MYSQL_USER", "go")
	MysqlPassword = GetenvOrDefault("MYSQL_PASSWORD", "Go1234_test")
	MysqlDataBase = GetenvOrDefault("MYSQL_DB", "go_test")
)

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

func NewNullInt(i int) NullInt {
	if i == 0 {
		return NullInt{}
	}
	return NullInt{
		Int:   i,
		Valid: true,
	}
}
