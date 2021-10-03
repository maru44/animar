package database

import (
	"animar/v1/internal/pkg/domain"

	"github.com/maru44/perr"
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
	if err != nil {
		return seasons, perr.Wrap(err, perr.InsufficientStorageWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var s domain.TSeasonRelation
		err := rows.Scan(
			&s.ID, &s.Year, &s.Season,
		)
		if err != nil {
			return seasons, perr.Wrap(err, perr.NotFound)
		}
		s.SeasonEng = domain.SeasonDictReverse[s.Season]
		seasons = append(seasons, s)
	}
	return
}
