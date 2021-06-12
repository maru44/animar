package tools

import (
	"animar/v1/configs"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

const slugLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func AccessDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s%s/%s", configs.MysqlUser, configs.MysqlPassword, configs.MysqlHost, configs.MysqlDataBase))
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
