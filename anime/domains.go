package anime

func ListAnimeDomain() []TAnime {
	rows := ListAnime()
	var animes []TAnime
	for rows.Next() {
		var ani TAnime
		err := rows.Scan(
			&ani.ID, &ani.Slug, &ani.Title, &ani.Content,
			&ani.OnAirState, &ani.SeriesId, &ani.Season,
			&ani.Stories, &ani.CreatedAt, &ani.UpdatedAt,
		)
		if err != nil {
			panic(err.Error())
		}
		animes = append(animes, ani)
	}

	defer rows.Close()

	return animes
}

/*
func ListAnimeWithUserWatchDomain() []TAnimeWithUserWatch {
	rows := ListAnime()
}
*/
