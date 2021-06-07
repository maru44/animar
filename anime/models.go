package anime

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TAnime struct {
	ID            int     `json:"id"`
	Slug          string  `json:"slug"`
	Title         string  `json:"title"`
	ThumbUrl      *string `json:"thumb_url,omitempty"`
	CopyRight     *string `json:"copyright,omitempty"`
	Abbreviation  *string `json:"abbreviation,omitempty"`
	Description   *string `json:"description,omitempty"`
	State         *string `json:"state,omitempty"`
	SeriesId      *int    `json:"series_id,omitempty"`
	CountEpisodes *int    `json:"count_episode,omitemptys"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
}

type TAnimeWithUserWatchReview struct {
	ID            int     `json:"id"`
	Slug          string  `json:"slug"`
	Title         string  `json:"title"`
	ThumbUrl      *string `json:"thumb_url,omitempty"`
	CopyRight     *string `json:"copyright,omitempty"`
	Abbreviation  *string `json:"abbreviation,omitempty"`
	Description   *string `json:"description,omitempty"`
	State         *string `json:"state,omitempty"`
	SeriesId      *int    `json:"series_id,omitempty"`
	CountEpisodes *int    `json:"count_episodes,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
	// watch from here
	Watch         int    `json:"watch,omitempty"`
	Star          int    `json:"star,omitempty"`
	ReviewContent string `json:"review_content,omitempty"`
}

type TAnimeAdmin struct {
	ID            int     `json:"id"`
	Slug          string  `json:"slug"`
	Title         string  `json:"title"`
	Abbreviation  *string `json:"abbreviation,omitempty"`
	Kana          *string `json:"kana,omitempty"`
	EngName       *string `json:"eng_name,omitempty"`
	ThumbUrl      *string `json:"thumb_url,omitempty"`
	CopyRight     *string `json:"copyright,omitempty"`
	Description   *string `json:"description,omitempty"`
	State         *string `json:"state,omitempty"`
	SeriesId      *int    `json:"series_id,omitempty"`
	CountEpisodes *int    `json:"count_episodes,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
}

type TAnimeMinimum struct {
	ID    int    `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

func ListAnime() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select id, slug, title, abbreviation, thumb_url, copyright, description, state, series_id, count_episodes, created_at, updated_at from animes")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func ListAnimeMinimum() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("SELECT id, slug, title from animes")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

// タイトル検索
func SearchAnime(title string) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query(
		"SELECT DISTINCT id, slug, title from animes where title like " +
			"'%" + title + "%' " +
			"OR kana like " +
			"'%" + title + "%' " +
			"OR eng_name like " +
			"'%" + title + "%' " +
			"limit 12",
	)

	// 1つでも NULLだと引っかからない
	// rows, err := db.Query(
	// 	"SELECT id,slug, title FROM animes where " +
	// 		"CONCAT(title,abbreviation,kana,eng_name) like" +
	// 		"'%" + title + "%'",
	// )

	// FULLTEXT indexでないとダメ
	// rows, err := db.Query(
	// 	"SELECT DISTINCT id, slug, title FROM animes " +
	// 		"WHERE MATCH (title,kana,eng_name,abbreviation) " +
	// 		"AGAINST ('+" + title + "')",
	// )
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func DetailAnime(id int) TAnime {
	db := tools.AccessDB()
	defer db.Close()

	var a TAnime
	err := db.QueryRow("SELECT id, slug, title, abbreviation, thumb_url, copyright, description, state, series_id, count_episodes, created_at, updated_at FROM animes WHERE id = ?", id).Scan(
		&a.ID, &a.Slug, &a.Title, &a.Abbreviation,
		&a.ThumbUrl, &a.CopyRight, &a.Description,
		&a.State, &a.SeriesId,
		&a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		a.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return a
}

func DetailAnimeBySlug(slug string) TAnime {
	db := tools.AccessDB()
	defer db.Close()

	var a TAnime
	err := db.QueryRow("SELECT id, slug, title, abbreviation, thumb_url, copyright, description, state, series_id, count_episodes, created_at, updated_at FROM animes WHERE slug = ?", slug).Scan(
		&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
		&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		a.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return a
}

// not be used to
func DetailAnimeBySlugWithUserWatchReview(slug string, userId string) TAnimeWithUserWatchReview {
	db := tools.AccessDB()
	defer db.Close()

	var a TAnimeWithUserWatchReview
	err := db.QueryRow(
		"SELECT anime.*, watch_states.watch, reviews.star, reviews.description AS review_content FROM animes "+
			"LEFT JOIN watch_states ON animes.id = watch_states.anime_id AND watch_states.user_id = ? "+
			"LEFT JOIN reviews ON animes.id = reviews.anime_id AND reviews.user_id = ? WHERE slug = ?",
		slug, userId, userId).Scan(
		&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.Description, &a.State, &a.CreatedAt, &a.UpdatedAt, &a.Watch, &a.Star, &a.ReviewContent,
	)

	switch {
	case err == sql.ErrNoRows:
		a.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return a
}

/************************************
             for admin
************************************/

func DetailAnimeAdmin(id int) TAnimeAdmin {
	db := tools.AccessDB()
	defer db.Close()

	var a TAnimeAdmin
	err := db.QueryRow("SELECT * FROM animes WHERE id = ?", id).Scan(
		&a.ID, &a.Slug, &a.Title, &a.Abbreviation,
		&a.Kana, &a.EngName, &a.ThumbUrl, &a.CopyRight, &a.Description,
		&a.State, &a.SeriesId,
		&a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		a.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return a
}

func InsertAnime(title string, abbreviation string, kana string, eng_name string, content string, State string, seriesId int, episodes int, copyright string, thumb_url string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO animes(title, slug, abbreviation, kana, eng_name, description, thumb_url, state, series_id, count_episodes, copyright) " +
			"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)
	defer stmtInsert.Close()

	slug := tools.GenRandSlug(12)
	exe, err := stmtInsert.Exec(
		title, slug, tools.NewNullString(abbreviation),
		tools.NewNullString(kana), tools.NewNullString(eng_name),
		tools.NewNullString(content), tools.NewNullString(thumb_url),
		tools.NewNullString(State), tools.NewNullInt(seriesId),
		tools.NewNullInt(episodes), tools.NewNullString(copyright),
	)

	insertedId, err := exe.LastInsertId()
	if err != nil {
		//panic(err.Error())
		fmt.Print(err)
	}
	return int(insertedId)
}

// @TODO insertにしたがってcolumn追加
func UpdateAnime(id int, title string, abbreviation string, kana string, eng_name string, content string, state string, seriesId int, episodes int, copyright string, thumbUrl string) int {
	db := tools.AccessDB()
	defer db.Close()

	exe, err := db.Exec(
		"UPDATE animes SET title = ?, abbreviation = ?, kana = ?, eng_name = ?, description = ?, thumb_url = ?, state = ?, series_id = ?, count_episodes = ?, copyright = ? WHERE id = ?",
		title, tools.NewNullString(abbreviation), tools.NewNullString(kana),
		tools.NewNullString(eng_name), tools.NewNullString(content),
		tools.NewNullString(thumbUrl), tools.NewNullString(state),
		tools.NewNullInt(seriesId), tools.NewNullInt(episodes),
		tools.NewNullString(copyright), id,
	)
	if err != nil {
		panic(err.Error())
	}
	updatedId, _ := exe.RowsAffected()
	return int(updatedId)
}

func DeleteAnime(id int) int {
	db := tools.AccessDB()
	defer db.Close()

	exe, err := db.Exec("DELETE FROM animes WHERE id = ?", id)
	if err != nil {
		panic(err.Error())
	}
	rowsAffect, _ := exe.RowsAffected()
	return int(rowsAffect)
}

func ListAnimeAdmin() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from animes")
	if err != nil {
		panic(err.Error())
	}
	return rows
}
