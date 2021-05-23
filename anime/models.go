package anime

import (
	"animar/v1/helper"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TAnime struct {
	ID          int
	Slug        string
	Title       string
	ThumbUrl    *string
	Abbribation *string
	Content     *string
	OnAirState  *int
	SeriesId    *int
	Season      *string
	Stories     *int
	CreatedAt   string
	UpdatedAt   string
}

type TAnimeWithUserWatchReview struct {
	ID          int
	Slug        string
	Title       string
	Abbribation *string
	ThumbUrl    *string
	Content     *string
	OnAirState  *int
	SeriesId    *int
	Season      *string
	Stories     *int
	CreatedAt   string
	UpdatedAt   string
	// watch from here
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
		&ani.ID, &ani.Slug, &ani.Title, &ani.Abbribation, &ani.ThumbUrl, &ani.Content,
		&ani.OnAirState, &ani.SeriesId, &ani.Season,
		&ani.Stories, &ani.CreatedAt, &ani.UpdatedAt,
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
		&ani.ID, &ani.Slug, &ani.Title, &ani.Abbribation, &ani.ThumbUrl, &ani.Content,
		&ani.OnAirState, &ani.SeriesId, &ani.Season,
		&ani.Stories, &ani.CreatedAt, &ani.UpdatedAt,
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
		&ani.ID, &ani.Slug, &ani.Title, &ani.Abbribation, &ani.Content, &ani.OnAirState, &ani.CreatedAt, &ani.UpdatedAt, &ani.Watch, &ani.Star, &ani.ReviewContent,
	)

	switch {
	case err == sql.ErrNoRows:
		ani.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return ani
}

func InsertAnime(title string, abbrevation string, content string, onAirState int, seriesId int, season string, stories int, thumb_url string) int {
	db := helper.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO anime(title, slug, abbrevation, content, on_air_state, series_id, season, stories) " +
			"VALUES(?, ?, ?, ?, ?, ?, ?)",
	)
	defer stmtInsert.Close()

	slug := helper.GenRandSlug(12)
	exe, err := stmtInsert.Exec(
		title, slug, helper.NewNullString(abbrevation).String,
		helper.NewNullString(thumb_url).String, helper.NewNullString(content).String,
		helper.NewNullInt(onAirState).Int, helper.NewNullInt(seriesId).Int,
		helper.NewNullString(season).String, helper.NewNullInt(stories).Int,
	)

	insertedId, err := exe.LastInsertId()
	if err != nil {
		//panic(err.Error())
		fmt.Print(err)
	}

	return int(insertedId)
}

// @TODO insertにしたがってcolumn追加
func UpdateAnime(slug string, abbrevation string, title string, content string, onAirState int, seriesId int, season string, stories int) int {
	db := helper.AccessDB()
	defer db.Close()

	exe, err := db.Exec(
		"UPDATE anime SET title = ?, abbrevation = ?, content = ?, on_air_state = ?, series_id = ?, season = ?, stories = ? WHERE slug = ?",
		title, helper.NewNullString(abbrevation).String, helper.NewNullString(content).String,
		helper.NewNullInt(onAirState).Int, helper.NewNullInt(seriesId).Int,
		helper.NewNullString(season).String, helper.NewNullInt(stories).Int, slug,
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
