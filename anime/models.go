package anime

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Anime struct {
	ID   int
	Name string
}

var MysqlUser = os.Getenv("MYSQL_USER")
var MysqlPassword = os.Getenv("MYSQL_PASSWORD")
var MysqlDataBase = os.Getenv("MYSQL_DB")

func AccessDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(http://localhost:3306)/%v", MysqlUser, MysqlPassword, MysqlDataBase))
	if err != nil {
		panic(err.Error())
	}

	return db
}
