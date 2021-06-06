package anime

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type TAnime struct {
	ID            int
	Slug          string
	Title         string
	ThumbUrl      *string
	Abbreviation  *string
	Description   *string
	State         *string
	SeriesId      *int
	CountEpisodes *int
	CreatedAt     string
	UpdatedAt     *string
}

type TAnimeWithUserWatchReview struct {
	ID            int
	Slug          string
	Title         string
	Abbreviation  *string
	ThumbUrl      *string
	Description   *string
	State         *string
	SeriesId      *int
	CountEpisodes *int
	CreatedAt     string
	UpdatedAt     *string
	// watch from here
	Watch         int
	Star          int
	ReviewContent string
}

type TAnimeAdmin struct {
	ID            int
	Slug          string
	Title         string
	Abbreviation  *string
	Kana          *string
	EngName       *string
	ThumbUrl      *string
	Description   *string
	State         *string
	SeriesId      *int
	CountEpisodes *int
	CreatedAt     string
	UpdatedAt     *string
}

type TAnimeMinimum struct {
	ID    int
	Slug  string
	Title string
}

func ListAnime() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select id, slug, title, abbreviation, thumb_url, description, state, series_id, count_episodes, created_at, updated_at from animes")
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
		"SELECT DISTINCT id, slug, title from animes where title like ?",
		"%"+title+"%",
		"OR kana like ?",
		"%"+title+"%",
		"OR eng_name like ?",
		"%"+title+"%",
		"limit 12",
	)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func DetailAnime(id int) TAnime {
	db := tools.AccessDB()
	defer db.Close()

	var ani TAnime
	err := db.QueryRow("SELECT id, slug, title, abbreviation, thumb_url, description, state, series_id, count_episodes, created_at, updated_at FROM animes WHERE id = ?", id).Scan(
		&ani.ID, &ani.Slug, &ani.Title, &ani.Abbreviation,
		&ani.ThumbUrl, &ani.Description,
		&ani.State, &ani.SeriesId,
		&ani.CountEpisodes, &ani.CreatedAt, &ani.UpdatedAt,
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
	db := tools.AccessDB()
	defer db.Close()

	var ani TAnime
	err := db.QueryRow("SELECT id, slug, title, abbreviation, thumb_url, description, on_air_state, series_id, count_episodes, created_at, updated_at FROM animes WHERE slug = ?", slug).Scan(
		&ani.ID, &ani.Slug, &ani.Title, &ani.Abbreviation, &ani.ThumbUrl, &ani.Description,
		&ani.State, &ani.SeriesId, &ani.CountEpisodes, &ani.CreatedAt, &ani.UpdatedAt,
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
	db := tools.AccessDB()
	defer db.Close()

	var ani TAnimeWithUserWatchReview
	err := db.QueryRow(
		"SELECT anime.*, watch_states.watch, tbl_reviews.star, tbl_reviews.description AS review_content FROM anime "+
			"LEFT JOIN watch_states ON animes.id = watch_states.anime_id AND watch_states.user_id = ? "+
			"LEFT JOIN tbl_reviews ON animes.id = tbl_reviews.anime_id AND tbl_reviews.user_id = ? WHERE slug = ?",
		slug, userId, userId).Scan(
		&ani.ID, &ani.Slug, &ani.Title, &ani.Abbreviation, &ani.Description, &ani.State, &ani.CreatedAt, &ani.UpdatedAt, &ani.Watch, &ani.Star, &ani.ReviewContent,
	)

	switch {
	case err == sql.ErrNoRows:
		ani.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return ani
}

/************************************
             for admin
************************************/

func DetailAnimeAdmin(id int) TAnimeAdmin {
	db := tools.AccessDB()
	defer db.Close()

	var ani TAnimeAdmin
	err := db.QueryRow("SELECT * FROM animes WHERE id = ?", id).Scan(
		&ani.ID, &ani.Slug, &ani.Title, &ani.Abbreviation,
		&ani.Kana, &ani.EngName, &ani.ThumbUrl, &ani.Description,
		&ani.State, &ani.SeriesId,
		&ani.CountEpisodes, &ani.CreatedAt, &ani.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		ani.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return ani
}

func InsertAnime(title string, abbreviation string, kana string, eng_name string, content string, State string, seriesId int, episodes int, thumb_url string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO animes(title, slug, abbreviation, kana, eng_name, content, thumb_url, state, series_id, count_episodes) " +
			"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)
	defer stmtInsert.Close()

	slug := tools.GenRandSlug(12)
	exe, err := stmtInsert.Exec(
		title, slug, tools.NewNullString(abbreviation),
		tools.NewNullString(kana), tools.NewNullString(eng_name),
		tools.NewNullString(content), tools.NewNullString(thumb_url),
		tools.NewNullString(State), tools.NewNullInt(seriesId),
		tools.NewNullInt(episodes),
	)

	insertedId, err := exe.LastInsertId()
	if err != nil {
		//panic(err.Error())
		fmt.Print(err)
	}
	return int(insertedId)
}

// @TODO insertにしたがってcolumn追加
func UpdateAnime(id int, title string, abbreviation string, kana string, eng_name string, content string, state string, seriesId int, episodes int, thumbUrl string) int {
	db := tools.AccessDB()
	defer db.Close()

	exe, err := db.Exec(
		"UPDATE animes SET title = ?, abbreviation = ?, kana = ?, eng_name = ?, content = ?, thumb_url = ?, state = ?, series_id = ?, count_episodes = ? WHERE id = ?",
		title, tools.NewNullString(abbreviation), tools.NewNullString(kana),
		tools.NewNullString(eng_name), tools.NewNullString(content),
		tools.NewNullString(thumbUrl), tools.NewNullString(state),
		tools.NewNullInt(seriesId),
		tools.NewNullInt(episodes), id,
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
