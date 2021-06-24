package connector

import (
	"animar/v1/configs"
	"database/sql"
	"fmt"
)

func AccessDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s%s/%s", configs.MysqlUser, configs.MysqlPassword, configs.MysqlHost, configs.MysqlDataBase))
	if err != nil {
		panic(err.Error())
	}
	return db
}
