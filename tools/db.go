package tools

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const slugLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

var (
	MysqlUser     = os.Getenv("MYSQL_USER")
	MysqlPassword = os.Getenv("MYSQL_PASSWORD")
	MysqlDataBase = os.Getenv("MYSQL_DB")
	MysqlHost     = os.Getenv("MYSQL_HOST")
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
