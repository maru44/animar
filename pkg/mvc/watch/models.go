package watch

import (
	"animar/v1/pkg/tools/connector"
	"animar/v1/pkg/tools/tools"
	"database/sql"
)

type TAudience struct {
	ID        int    `json:"id"`
	State     int    `json:"state"`
	AnimeId   int    `json:"anime_id"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TAudienceCount struct {
	State int `json:"state"`
	Count int `json:"count"`
}

type TAudienceJoinAnime struct {
	ID        int     `json:"id"`
	State     int     `json:"state"`
	AnimeId   int     `json:"anime_id"`
	UserId    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Title     string  `json:"title"`
	Slug      string  `json:"slug"`
	Content   *string `json:"content"`
	AState    *string `json:"anime_state"`
}

// Count group by animeId
// fiter by animeId
func AnimeWatchCounts(animeId int) *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select state, count(state) from audiences WHERE anime_id = ? GROUP BY state", animeId)
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

// List
// filter by user
func OnesAnimeWatchList(userId string) *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from audiences WHERE user_id = ?", userId)
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

func OnesAnimeWatchJoinList(userId string) *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()
	rows, err := db.Query(
		"SELECT audiences.*, animes.title, animes.slug, animes.description, animes.state "+
			"FROM audiences LEFT JOIN animes ON audiences.anime_id = animes.id WHERE user_id = ?", userId,
	)
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

func WatchDetail(userId string, animeId int) int {
	db := connector.AccessDB()
	defer db.Close()

	var w TAudience
	err := db.QueryRow("Select * from audiences WHERE user_id = ? AND anime_id = ?", userId, animeId).Scan(
		&w.ID, &w.State, &w.AnimeId, &w.UserId, &w.CreatedAt, &w.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		return -1
	case err != nil:
		tools.ErrorLog(err)
		return 0
	default:
		return w.State
	}
}

// Post
func InsertWatch(animeId int, state int, userId string) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO audiences(state, anime_id, user_id) VALUES(?, ?, ?)")
	defer stmt.Close()

	exe, err := stmt.Exec(state, animeId, userId)

	_, err = exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
	}

	return state
}

// create or update
func UpsertWatch(animeId int, state int, userId string) int {
	db := connector.AccessDB()
	defer db.Close()

	var w TAudience
	err := db.QueryRow("Select * from audiences WHERE user_id = ? AND anime_id = ?", userId, animeId).Scan(
		&w.ID, &w.State, &w.AnimeId, &w.UserId, &w.CreatedAt, &w.UpdatedAt,
	)
	switch {
	case err == sql.ErrNoRows:
		// create
		return InsertWatch(animeId, state, userId)
	case err != nil:
		tools.ErrorLog(err)
		return -1
	default:
		// update
		stmt, err := db.Prepare("UPDATE audiences SET state = ? WHERE user_id = ? AND anime_id = ?")
		defer stmt.Close()
		_, err = stmt.Exec(state, userId, animeId)
		if err != nil {
			tools.ErrorLog(err)
		}
		return state
	}
}

func DeleteWatch(animeId int, userId string) bool {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM audiences WHERE anime_id = ? AND user_id = ?")
	_, err = stmt.Exec(animeId, userId)
	if err != nil {
		//panic(err.Error())
		return false
	}
	return true
}
