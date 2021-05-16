package watch

import (
	"animar/v1/helper"
	"database/sql"
	"fmt"
	"reflect"
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

func WatchDetail(userId string, animeId int) int {
	db := helper.AccessDB()
	defer db.Close()

	var watch TWatch
	fmt.Println(userId, reflect.TypeOf(animeId))
	err := db.QueryRow("Select * from watch_states WHERE user_id = ? AND anime_id = ?", userId, animeId).Scan(
		&watch.ID, &watch.Watch, &watch.AnimeId, &watch.UserId, &watch.CreatedAt, &watch.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		return -1
	case err != nil:
		panic(err.Error())
	default:
		return watch.Watch
	}
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

	return watch
}

// create or update
func UpsertWatch(animeId int, watch int, userId string) int {
	db := helper.AccessDB()
	defer db.Close()

	var w TWatch
	err := db.QueryRow("Select * from watch_states WHERE user_id = ? AND anime_id = ?", userId, animeId).Scan(
		&w.ID, &w.Watch, &w.AnimeId, &w.UserId, &w.CreatedAt, &w.UpdatedAt,
	)
	switch {
	case err == sql.ErrNoRows:
		// create
		return InsertWatch(animeId, watch, userId)
	case err != nil:
		panic(err.Error())
	default:
		// update
		stmtUpdate, err := db.Prepare("UPDATE watch_states SET watch = ? WHERE user_id = ? AND anime_id = ?")
		defer stmtUpdate.Close()
		exe, err := stmtUpdate.Exec(watch, userId, animeId)
		if err != nil {
			panic(err.Error())
		}
		fmt.Print(exe)
		return watch
	}
}

func DeleteWatch(animeId int, userId string) bool {
	db := helper.AccessDB()
	defer db.Close()

	exe, err := db.Exec("DELETE FROM watch_states WHERE anime_id = ? AND user_id = ?", animeId, userId)
	if err != nil {
		//panic(err.Error())
		return false
	}
	fmt.Print(exe.RowsAffected())
	return true
}
