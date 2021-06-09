package tools

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

const slugLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

var (
	MysqlUser     = GetenvOrDefault("MYSQL_USER", "go")
	MysqlPassword = GetenvOrDefault("MYSQL_PASSWORD", "Go1234_test")
	MysqlDataBase = GetenvOrDefault("MYSQL_DB", "go_test")
	MysqlHost     = GetenvOrDefault("MYSQL_HOST", "@tcp(127.0.0.1:3306)")
)

func AccessDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s%s/%s", MysqlUser, MysqlPassword, MysqlHost, MysqlDataBase))
	if err != nil {
		panic(err.Error())
	}
	return db
}

// insert & update用
func NewNullString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// insert & update用
func NewNullInt(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func GenRandSlug(n int) string {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = slugLetters[rand.Intn(len(slugLetters))]
	}
	return string(b)
}
