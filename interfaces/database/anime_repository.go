package database

import (
	"animar/v1/domain"
	"animar/v1/tools/tools"
)

type AnimeRepository struct {
	SqlHandler
}

func (repo *AnimeRepository) ListAll() (animes domain.TAnimes, err error) {
	rows, err := repo.Query(
		"SELECT id, slug, title, abbreviation, thumb_url, copyright, description, state, series_id, " +
			"count_episodes, created_at, updated_at FROM animes ORDER BY created_at ASC",
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var a domain.TAnime
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListMinimumAll() (animes domain.TAnimeMinimums, err error) {
	rows, err := repo.Query(
		"SELECT id, slug, title FROM animes",
	)
	defer rows.Close()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var a domain.TAnimeMinimum
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListOnAirAll() (animes domain.TAnimes, err error) {
	rows, err := repo.Query(
		"SELECT id, slug, title, abbreviation, thumb_url, copyright, " +
			"description, state, series_id, count_episodes, created_at, updated_at " +
			"FROM animes WHERE state = 'now' ORDER BY created_at ASC",
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var a domain.TAnime
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListMinimumSearch(title string) (animes domain.TAnimeMinimums, err error) {
	rows, err := repo.Query(
		"SELECT DISTINCT id, slug, title from animes where title like " +
			"'%" + title + "%' " +
			"OR kana like " +
			"'%" + title + "%' " +
			"OR eng_name like " +
			"'%" + title + "%' " +
			"limit 12",
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var a domain.TAnimeMinimum
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListSearch(title string) (animes domain.TAnimes, err error) {
	rows, err := repo.Query(
		"SELECT DISTINCT id, slug, title from animes where title like " +
			"'%" + title + "%' " +
			"OR kana like " +
			"'%" + title + "%' " +
			"OR eng_name like " +
			"'%" + title + "%' " +
			"limit 12",
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var a domain.TAnime
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListBySeason(year string, season string) (animes domain.TAnimes, err error) {
	seasonJp, _ := domain.SeasonDict[season]
	rows, err := repo.Query(
		"SELECT animes.id as id, slug, title, abbreviation, thumb_url, copyright, "+
			"description, state, series_id, count_episodes, "+
			"animes.created_at as created_at, animes.updated_at as updated_at "+
			"FROM seasons "+
			"LEFT JOIN relation_anime_season as rel ON seasons.id = rel.season_id "+
			"LEFT JOIN animes ON rel.anime_id = animes.id "+
			"WHERE seasons.year = ? AND seasons.season = ?", year, seasonJp,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var a domain.TAnime
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		animes = append(animes, a)
	}
	return
}

/*******************************
             Detail
*******************************/

func (repo *AnimeRepository) FindById(id int) (a domain.TAnime, err error) {
	row, err := repo.QueryRow(
		"SELECT id, slug, title, abbreviation, thumb_url, copyright, description, state, series_id, "+
			"count_episodes, created_at, updated_at FROM animes WHERE id = ?",
		id,
	)
	if err != nil {
		tools.ErrorLog(err)
	}
	row.Scan(
		&a.ID, &a.Slug, &a.Title, &a.Abbreviation,
		&a.ThumbUrl, &a.CopyRight, &a.Description,
		&a.State, &a.SeriesId,
		&a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
	)
	return
}

func (repo *AnimeRepository) FindBySlug(slug string) (a domain.TAnimeWithSeries, err error) {
	row, err := repo.QueryRow(
		"SELECT animes.id as id, slug, title, abbreviation, thumb_url, copyright, description, state, series_id, count_episodes, animes.created_at, animes.updated_at, "+
			"series_name FROM animes "+
			"LEFT JOIN series on animes.series_id = series.id "+
			"WHERE slug = ?", slug,
	)
	if err != nil {
		tools.ErrorLog(err)
	}
	row.Scan(
		&a.ID, &a.Slug, &a.Title, &a.Abbreviation,
		&a.ThumbUrl, &a.CopyRight, &a.Description,
		&a.State, &a.SeriesId, &a.CountEpisodes,
		&a.CreatedAt, &a.UpdatedAt, &a.SeriesName,
	)
	return
}
