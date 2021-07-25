package database

import (
	"animar/v1/internal/pkg/domain"
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
		domain.LogWriter(err.Error())
		return
	}
	for rows.Next() {
		var s domain.TSeasonRelation
		err := rows.Scan(
			&s.ID, &s.Year, &s.Season,
		)
		s.SeasonEng = domain.SeasonDictReverse[s.Season]
		if err != nil {
			domain.LogWriter(err.Error())
		}
		seasons = append(seasons, s)
	}
	return
}
