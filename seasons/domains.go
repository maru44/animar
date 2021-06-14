package seasons

import "animar/v1/tools/tools"

func ListSeasonDomain() []TSeason {
	rows := listSeason()
	var seasons []TSeason
	for rows.Next() {
		var s TSeason
		err := rows.Scan(
			&s.ID, &s.Year, &s.Season, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		seasons = append(seasons, s)
	}

	defer rows.Close()
	return seasons
}

/************************************
             relation
************************************/

func SeasonByAnimeIdDomain(animeId int) []TSeasonRelation {
	rows := relationSeasonByAnime(animeId)
	var seasons []TSeasonRelation
	for rows.Next() {
		var s TSeasonRelation
		err := rows.Scan(
			&s.ID, &s.Year, &s.Season,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		seasons = append(seasons, s)
	}

	defer rows.Close()
	return seasons
}
