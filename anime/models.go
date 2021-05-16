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

type TAnimeWithUserWatchReview struct {
	ID            int
	Slug          string
	Title         string
	Content       *string
	OnAirState    *int
	CreatedAt     string
	UpdatedAt     string
	Watch         int
	Star          int
	ReviewContent string
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

// not be used to
func DetailAnimeBySlugWithUserWatchReview(slug string, userId string) TAnimeWithUserWatchReview {
	db := helper.AccessDB()
	defer db.Close()

	var ani TAnimeWithUserWatchReview
	err := db.QueryRow(
		"SELECT anime.*, watch_states.watch, tbl_reviews.star, tbl_reviews.content AS review_content FROM anime "+
			"LEFT JOIN watch_states ON anime.id = watch_states.anime_id AND watch_states.user_id = ? "+
			"LEFT JOIN tbl_reviews ON anime.id = tbl_reviews.anime_id AND tbl_reviews.user_id = ? WHERE slug = ?",
		slug, userId, userId).Scan(
		&ani.ID, &ani.Slug, &ani.Title, &ani.Content, &ani.OnAirState, &ani.CreatedAt, &ani.UpdatedAt, &ani.Watch, &ani.Star, &ani.ReviewContent,
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
		title, slug, helper.NewNullString(content).String, helper.NewNullInt(onAirState).Int,
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
		title, helper.NewNullString(content).String, slug, helper.NewNullInt(onAirState).Int,
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
