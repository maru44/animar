package anime

import (
	"animar/v1/helper"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TAnime struct {
	ID        int
	Title     string
	Content   string
	CreatedAt string
}

var MysqlUser = helper.GetenvOrDefault("MYSQL_USER", "go")
var MysqlPassword = helper.GetenvOrDefault("MYSQL_PASSWORD", "Go1234_test")
var MysqlDataBase = helper.GetenvOrDefault("MYSQL_DB", "go_test")

func AccessDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", MysqlUser, MysqlPassword, MysqlDataBase))
	if err != nil {
		panic(err.Error())
	}

	return db
}

func ListAnime() *sql.Rows {
	db := AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from anime")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func DetailAnime(id int) TAnime {
	db := AccessDB()
	defer db.Close()

	var ani TAnime
	nullContent := new(sql.NullString)
	err := db.QueryRow("SELECT * FROM anime WHERE id = ?", id).Scan(&ani.ID, &ani.Title, nullContent, &ani.CreatedAt)

	switch {
	case err == sql.ErrNoRows:
		ani.ID = 0
	case err != nil:
		panic(err.Error())
	default:
		ani.Content = nullContent.String
	}
	return ani
}

func InsertAnime(title string, content string) sql.Result {
	db := AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO anime (title, content) VALUES (:title, :content)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(title, content)
	if err != nil {
		panic(err.Error())
	}
	return exe
}
