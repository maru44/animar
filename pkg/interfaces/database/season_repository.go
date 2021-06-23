package database

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/tools/tools"
)

type SeasonRepository struct {
	SqlHandler
}

func (repo *SeasonRepository) FilterByAnimeId(animeId int) (seasons []domain.TSeasonRelation, err error) {
	rows, err := repo.Query(
		"SELECT seasons.id, seasons.year, seasons.season FROM relation_anime_season "+
			"LEFT JOIN seasons ON relation_anime_season.season_id = seasons.id "+
			"WHERE anime_id = ?", animeId,
	)
	defer rows.Close()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var s domain.TSeasonRelation
		err := rows.Scan(
			&s.ID, &s.Year, &s.Season,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		seasons = append(seasons, s)
	}
	return
}
