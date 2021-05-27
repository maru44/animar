package anime

func ListAnimeDomain() []TAnime {
	rows := ListAnime()
	var animes []TAnime
	for rows.Next() {
		var ani TAnime
		err := rows.Scan(
			&ani.ID, &ani.Slug, &ani.Title, &ani.Abbribation, &ani.ThumbUrl, &ani.Content,
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

func ListAnimeMinimumDomain() []TAnimeMinimum {
	rows := ListAnimeMinimum()
	var animes []TAnimeMinimum
	for rows.Next() {
		var ani TAnimeMinimum
		err := rows.Scan(
			&ani.ID, &ani.Title,
		)
		if err != nil {
			panic(err.Error())
		}
		animes = append(animes, ani)
	}

	defer rows.Close()
	return animes
}
