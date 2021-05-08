package anime

import (
	"animar/v1/helper"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type TAnime struct {
	ID        int
	Title     string
	Content   string
	CreatedAt string
	UpdatedAt string
}

func ListAnime() *sql.Rows {
	db := helper.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from anime")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func DetailAnime(id int) TAnime {
	db := helper.AccessDB()
	defer db.Close()

	var ani TAnime
	nullContent := new(sql.NullString)
	err := db.QueryRow("SELECT * FROM anime WHERE id = ?", id).Scan(&ani.ID, &ani.Title, nullContent, &ani.CreatedAt, &ani.UpdatedAt)

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

func InsertAnime(title string, content string) int {
	db := helper.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO anime(title, content) VALUES(?, ?)")
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(title, helper.NewNullString(content))

	insertedId, err := exe.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	return int(insertedId)
}
