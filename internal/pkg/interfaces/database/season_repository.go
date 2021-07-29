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
		lg := domain.NewErrorLog(err.Error(), "")
		lg.Logging()
		return
	}
	for rows.Next() {
		var s domain.TSeasonRelation
		err := rows.Scan(
			&s.ID, &s.Year, &s.Season,
		)
		s.SeasonEng = domain.SeasonDictReverse[s.Season]
		if err != nil {
			lg := domain.NewErrorLog(err.Error(), "")
			lg.Logging()
		}
		seasons = append(seasons, s)
	}
	return
}
