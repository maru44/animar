package anime

import "animar/v1/tools"

func ListAnimeDomain() []TAnime {
	rows := ListAnime()
	var animes []TAnime
	for rows.Next() {
		var a TAnime
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, a)
	}

	defer rows.Close()
	return animes
}

func ListAnimeMinimumDomain() []TAnimeMinimum {
	rows := ListAnimeMinimum()
	var animes []TAnimeMinimum
	for rows.Next() {
		var a TAnimeMinimum
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, a)
	}

	defer rows.Close()
	return animes
}

func ListAnimeMinimumDomainByTitle(title string) []TAnimeMinimum {
	rows := SearchAnime(title)
	var animes []TAnimeMinimum
	for rows.Next() {
		var a TAnimeMinimum
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, a)
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
			&ani.Kana, &ani.EngName, &ani.ThumbUrl, &ani.Description,
			&ani.State, &ani.SeriesId,
			&ani.CountEpisodes, &ani.CreatedAt, &ani.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, ani)
	}

	defer rows.Close()
	return animes
}
