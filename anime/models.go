package anime

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Anime struct {
	ID int
	Name string
}

const MysqlUser = os.Getenv("MYSQL_USER")
const MysqlPassword = os.Getenv("MYSQL_PASSWORD")
const MysqlDataBase = os.Getenv("MYSQL_DB")

func AccessDB() {
	db, err := sql.Open("mysql", fmt.Fprintf("%v:%v@tcp(127.0.0.1:3306)/%v"), MysqlUser, MysqlPassword, MysqlDataBase)
	if err != nil {
		panic(err.Error())
	}
}

