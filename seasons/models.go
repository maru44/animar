package seasons

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"
)

type TSeasonRelation struct {
	ID     int    `json:"id"`
	Year   string `json:"year"`
	Season string `json:"season"`
}

type TSeason struct {
	ID        int    `json:"id"`
	Year      string `json:"year"`
	Season    string `json:"season"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TSeasonInput struct {
	Year   string `json:"year"`
	Season string `json:"season"`
}

func InsertSeason(year string, season string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare(
		"INSERT INTO seasons(year, season) VALUES(?, ?)",
	)
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(
		year, season,
	)
	insertedId, err := exe.LastInsertId()
	if err != nil {
		fmt.Print(err)
	}
	return int(insertedId)
}

func ListSeason() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from seasons")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func DetailSeason(id int) TSeason {
	db := tools.AccessDB()
	defer db.Close()

	var s TSeason
	err := db.QueryRow(
		"SELECT * FROM seasons WHERE id = ?", id,
	).Scan(
		s.ID, s.Year, s.Season, s.CreatedAt, s.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		s.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return s
}

/************************************
             relation
************************************/

func RelationSeasonByAnime(animeId int) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query(
		"SELECT seasons.id, seasons.year, seasons.season FROM relation_anime_season "+
			"LEFT JOIN seasons ON relation_anime_season.season_id = seasons.id"+
			"WHERE anime_id = ?", animeId,
	)
	if err != nil {
		panic(err.Error())
	}
	return rows
}
