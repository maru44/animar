package database

import (
	"animar/v1/internal/pkg/domain"
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
		return
	}
	for rows.Next() {
		var a domain.TAnime
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return animes, err
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
	for rows.Next() {
		var a domain.TAnime
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return animes, err
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListMinimumSearch(title string) (animes domain.TAnimeMinimums, err error) {
	rows, err := repo.Query(
		"SELECT DISTINCT id, slug, title FROM animes "+
			"WHERE title LIKE "+
			"CONCAT('%', ?, '%') "+
			"OR kana LIKE "+
			"CONCAT('%', ?, '%') "+
			"OR eng_name LIKE "+
			"CONCAT('%', ?, '%') "+
			"LIMIT 12",
		title, title, title,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		var a domain.TAnimeMinimum
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title,
		)
		if err != nil {
			return animes, err
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListSearch(title string) (animes domain.TAnimes, err error) {
	rows, err := repo.Query(
		"SELECT DISTINCT id, slug, title, abbreviation, thumb_url, copyright, "+
			"description, state, series_id, count_episodes, created_at, updated_at "+
			"FROM animes "+
			"WHERE title LIKE "+
			"CONCAT('%', ?, '%') "+
			"OR kana LIKE "+
			"CONCAT('%', ?, '%') "+
			"OR eng_name LIKE "+
			"CONCAT('%', ?, '%') "+
			"ORDER BY created_at DESC",
		title, title, title,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		var a domain.TAnime
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return animes, err
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListBySeason(year string, season string) (animes domain.TAnimes, err error) {
	seasonJp, _ := domain.SeasonDict[season]
	rows, err := repo.SqlHandler.Query(
		"SELECT animes.id as id, slug, title, abbreviation, thumb_url, copyright, "+
			"description, state, series_id, count_episodes, "+
			"animes.created_at as created_at, animes.updated_at as updated_at "+
			"FROM seasons "+
			"LEFT JOIN relation_anime_season as rel ON seasons.id = rel.season_id "+
			"LEFT JOIN animes ON rel.anime_id = animes.id "+
			"WHERE seasons.year = ? AND seasons.season = ?", year, seasonJp,
	)
	if err != nil {
		return
	}
	for rows.Next() {
		var a domain.TAnime
		err := rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation, &a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return animes, err
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListBySeries(id int) (animes []domain.TAnimeWithSeries, err error) {
	rows, err := repo.Query(
		"SELECT animes.id as id, slug, title, abbreviation, thumb_url, copyright, "+
			"description, state, series_id, count_episodes, "+
			"animes.created_at as created_at, animes.updated_at as updated_at, series.series_name "+
			"FROM series "+
			"LEFT JOIN animes ON series.id = animes.series_id "+
			"WHERE series.id = ?", id,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var a domain.TAnimeWithSeries
		err = rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation,
			&a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes,
			&a.CreatedAt, &a.UpdatedAt, &a.SeriesName,
		)
		if err != nil {
			return
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListByCompany(engName string) (animes domain.TAnimes, err error) {
	rows, err := repo.Query(
		"SELECT animes.id as id, slug, title, abbreviation, thumb_url, copyright, "+
			"description, state, series_id, count_episodes, "+
			"animes.created_at as created_at, animes.updated_at as updated_at "+
			"FROM companies AS co "+
			"LEFT JOIN animes ON co.id = animes.company_id "+
			"WHERE co.eng_name = ?", engName,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var a domain.TAnime
		err = rows.Scan(
			&a.ID, &a.Slug, &a.Title, &a.Abbreviation,
			&a.ThumbUrl, &a.CopyRight, &a.Description,
			&a.State, &a.SeriesId, &a.CountEpisodes,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return
		}
		animes = append(animes, a)
	}
	return
}

func (repo *AnimeRepository) ListMinimum() (animes domain.TAnimeMinimums, err error) {
	rows, err := repo.Query(
		"SELECT id, slug, title from animes",
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var a domain.TAnimeMinimum
		err = rows.Scan(&a.ID, &a.Slug, &a.Title)
		if err != nil {
			return
		}
		animes = append(animes, a)
	}
	return
}

/*******************************
             Detail
*******************************/

func (repo *AnimeRepository) FindById(id int) (a domain.TAnime, err error) {
	rows, err := repo.Query(
		"SELECT id, slug, title, abbreviation, thumb_url, copyright, description, state, series_id, "+
			"count_episodes, created_at, updated_at FROM animes WHERE id = ?",
		id,
	)
	if err != nil {
		return a, domain.NewWrapError(err, domain.MySqlConnectionError)
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(
		&a.ID, &a.Slug, &a.Title, &a.Abbreviation,
		&a.ThumbUrl, &a.CopyRight, &a.Description,
		&a.State, &a.SeriesId,
		&a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return a, domain.NewWrapError(err, domain.DataNotFoundError)
	}
	return
}

func (repo *AnimeRepository) FindBySlug(slug string) (a domain.TAnimeWithSeries, err error) {
	rows, err := repo.Query(
		"SELECT a.id as id, slug, title, abbreviation, thumb_url, copyright, description, state, series_id, "+
			"count_episodes, hash_tag, twitter_url, a.official_url, a.created_at, a.updated_at, "+
			"series_name, co.name, co.eng_name "+
			"FROM animes AS a "+
			"LEFT JOIN series on a.series_id = series.id "+
			"LEFT JOIN companies AS co ON a.company_id = co.id "+
			"WHERE slug = ?", slug,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(
		&a.ID, &a.Slug, &a.Title, &a.Abbreviation,
		&a.ThumbUrl, &a.CopyRight, &a.Description,
		&a.State, &a.SeriesId, &a.CountEpisodes,
		&a.HashTag, &a.TwitterUrl, &a.OfficialUrl,
		&a.CreatedAt, &a.UpdatedAt, &a.SeriesName,
		&a.CompanyName, &a.CompanyEngName,
	)
	if err != nil {
		return
	}
	return
}

/***********************
          review
***********************/

func (repo *AnimeRepository) ReviewFilterByAnime(animeId int, userId string) (reviews domain.TReviews, err error) {
	rows, err := repo.Query(
		"Select * from reviews WHERE anime_id = ? AND (user_id != ? OR user_id IS NULL)",
		animeId, userId,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var r domain.TReview
		err = rows.Scan(&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return
		}
		reviews = append(reviews, r)
	}
	return
}
