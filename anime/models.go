package anime

import (
	"animar/v1/helper"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type TAnime struct {
	ID         int
	Slug       string
	Title      string
	Content    *string
	OnAirState *int
	CreatedAt  string
	UpdatedAt  string
}

type TAnimeWithUserWatch struct {
	ID         int
	Slug       string
	Title      string
	Content    *string
	OnAirState *int
	CreatedAt  string
	UpdatedAt  string
	Watch      int
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
	err := db.QueryRow("SELECT * FROM anime WHERE id = ?", id).Scan(
		&ani.ID, &ani.Slug, &ani.Title, &ani.Content, &ani.OnAirState, &ani.CreatedAt, &ani.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		ani.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return ani
}

func DetailAnimeBySlug(slug string) TAnime {
	db := helper.AccessDB()
	defer db.Close()

	var ani TAnime
	err := db.QueryRow("SELECT * FROM anime WHERE slug = ?", slug).Scan(
		&ani.ID, &ani.Slug, &ani.Title, &ani.Content, &ani.OnAirState, &ani.CreatedAt, &ani.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		ani.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return ani
}

func DetailAnimeBySlugWithUserWatch(slug string, userId string) TAnimeWithUserWatch {
	db := helper.AccessDB()
	defer db.Close()

	var ani TAnimeWithUserWatch
	err := db.QueryRow("SELECT * FROM anime AS T1 WHERE slug = ? LEFT JOIN watch_states AS T2 ON T1.id = T2.anime_id AND T2.user_id = ?", slug, userId).Scan(
		&ani.ID, &ani.Slug, &ani.Title, &ani.Content, &ani.OnAirState, &ani.CreatedAt, &ani.UpdatedAt, &ani.Watch,
	)

	switch {
	case err == sql.ErrNoRows:
		ani.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return ani
}

func InsertAnime(title string, content string, onAirState int) int {
	db := helper.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO anime(title, slug, content) VALUES(?, ?, ?)")
	defer stmtInsert.Close()

	slug := helper.GenRandSlug(12)
	exe, err := stmtInsert.Exec(
		title, slug, helper.NewNullString(content), helper.NewNullInt(onAirState),
	)

	insertedId, err := exe.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	return int(insertedId)
}

func UpdateAnime(slug string, title string, content string, onAirState int) int {
	db := helper.AccessDB()
	defer db.Close()

	exe, err := db.Exec(
		"UPDATE anime SET title = ?, content = ?, on_air_state = ? WHERE slug = ?",
		title, content, slug, onAirState,
	)
	if err != nil {
		panic(err.Error())
	}
	updatedId, _ := exe.RowsAffected()
	return int(updatedId)
}

func DeleteAnime(id int) int {
	db := helper.AccessDB()
	defer db.Close()

	exe, err := db.Exec("DELETE FROM anime WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}
