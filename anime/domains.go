package anime

func ListAnimeDomain() []TAnime {
	rows := ListAnime()
	var animes []TAnime
	for rows.Next() {
		var ani TAnime
		err := rows.Scan(
			&ani.ID, &ani.Slug, &ani.Title, &ani.Abbreviation, &ani.ThumbUrl, &ani.Content,
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
			&ani.ID, &ani.Slug, &ani.Title,
		)
		if err != nil {
			panic(err.Error())
		}
		animes = append(animes, ani)
	}

	defer rows.Close()
	return animes
}

func ListAnimeMinimumDomainByTitle(title string) []TAnimeMinimum {
	rows := SearchAnime(title)
	var animes []TAnimeMinimum
	for rows.Next() {
		var ani TAnimeMinimum
		err := rows.Scan(
			&ani.ID, &ani.Slug, &ani.Title,
		)
		if err != nil {
			panic(err.Error())
		}
		animes = append(animes, ani)
	}

	defer rows.Close()
	return animes
}

/************************************
             for admin
************************************/

func ListAnimeAdminDomain() []TAnimeAdmin {
	rows := ListAnimeAdmin()
	var animes []TAnimeAdmin
	for rows.Next() {
		var ani TAnimeAdmin
		err := rows.Scan(
			&ani.ID, &ani.Slug, &ani.Title, &ani.Abbreviation,
			&ani.Kana, &ani.EngName, &ani.ThumbUrl, &ani.Content,
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
