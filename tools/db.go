package tools

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

const slugLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

type NullInt struct {
	Int   *int
	Valid bool
}

type NullStringMine struct {
	String *string
	Valid  bool
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

// insert & update用
func NewNullString(s string) NullStringMine {
	if s == "" {
		return NullStringMine{
			String: nil,
			Valid:  false,
		}
	}
	return NullStringMine{
		String: &s,
		Valid:  true,
	}
}

// insert & update用
func NewNullInt(i int) NullInt {
	if i == 0 {
		return NullInt{
			Int:   nil,
			Valid: false,
		}
	}
	return NullInt{
		Int:   &i,
		Valid: true,
	}
}

func GenRandSlug(n int) string {
	b := make([]byte, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = slugLetters[rand.Intn(len(slugLetters))]
	}
	return string(b)
}
