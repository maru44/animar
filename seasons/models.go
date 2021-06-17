package seasons

import (
	"animar/v1/tools/connector"
	"animar/v1/tools/tools"
	"database/sql"
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

var SeasonDict = map[string]string{
	"spring": "春",
	"summer": "夏",
	"fall":   "秋",
	"winter": "冬",
}

func insertSeason(year string, season string) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare(
		"INSERT INTO seasons(year, season) VALUES(?, ?)",
	)
	defer stmt.Close()

	exe, err := stmt.Exec(
		year, season,
	)
	insertedId, err := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
	}
	return int(insertedId)
}

func listSeason() *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from seasons")
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

func detailSeason(id int) TSeason {
	db := connector.AccessDB()
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
		tools.ErrorLog(err)
	}
	return s
}

/************************************
             relation
************************************/

func relationSeasonByAnime(animeId int) *sql.Rows {
	db := connector.AccessDB()
	defer db.Close()
	rows, err := db.Query(
		"SELECT seasons.id, seasons.year, seasons.season FROM relation_anime_season "+
			"LEFT JOIN seasons ON relation_anime_season.season_id = seasons.id "+
			"WHERE anime_id = ?", animeId,
	)
	if err != nil {
		tools.ErrorLog(err)
	}
	return rows
}

func insertRelation(seasonId int, animeId int) int {
	db := connector.AccessDB()
	defer db.Close()

	stmt, err := db.Prepare(
		"INSERT INTO relation_anime_season(season_id, anime_id) VALUES(?, ?)",
	)
	defer stmt.Close()

	exe, err := stmt.Exec(
		seasonId, animeId,
	)
	insertedId, _ := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
	}
	return int(insertedId)
}
