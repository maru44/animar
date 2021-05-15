package watch

import (
	"animar/v1/helper"
	"database/sql"
	"fmt"
)

type TWatch struct {
	ID        int
	Watch     int
	AnimeId   int
	UserId    string
	CreatedAt string
	UpdatedAt string
}

type TWatchCount struct {
	State int // watch state
	Count int // count of watcher
}

// Count group by animeId
// fiter by animeId
func AnimeWatchCounts(animeId int) *sql.Rows {
	db := helper.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select watch, count(watch) from watch_states WHERE anime_id = ? GROUP BY WATCH", animeId)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

// List
// filter by user
func OnesAnimeWatchList(userId string) *sql.Rows {
	db := helper.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from watch_states WHERE user_id = ?", userId)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func WatchDetail(userId string, animeId int) *TWatch {
	db := helper.AccessDB()
	defer db.Close()

	var watch TWatch
	err := db.QueryRow("Select * from watch_states WHERE user_id = ?, animeId = ?", userId, animeId).Scan(&watch.ID, &watch.Watch, &watch.AnimeId, &watch.UserId, &watch.CreatedAt, &watch.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		panic(err.Error())
	}
	return &watch
}

// Post
func InsertWatch(animeId int, watch int, userId string) int {
	db := helper.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO watch_states(watch, anime_id, user_id) VALUES(?, ?, ?)")
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(watch, animeId, userId)

	insertedId, err := exe.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	fmt.Print(insertedId)

	//return int(insertedId)
	return watch
}
