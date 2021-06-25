package database

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/tools/tools"
)

type AdminAnimeRepository struct {
	SqlHandler
}

type AdminPlatformRepository struct {
	SqlHandler
}

type AdminSeasonRepository struct {
	SqlHandler
}

type AdminSeriesRepository struct {
	SqlHandler
}

/************************
         anime
*************************/

func (repo *AdminAnimeRepository) ListAll() (animes domain.TAnimes, err error) {
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

func (repo *AdminAnimeRepository) FindById(id int) (a domain.TAnimeAdmin, err error) {
	rows, err := repo.Query(
		"SELECT * FROM animes WHERE id = ?", id,
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rows.Next()
	err = rows.Scan(
		&a.ID, &a.Slug, &a.Title, &a.Abbreviation,
		&a.Kana, &a.EngName, &a.ThumbUrl, &a.CopyRight, &a.Description,
		&a.State, &a.SeriesId,
		&a.CountEpisodes, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	return
}

func (repo *AdminAnimeRepository) Insert(a domain.TAnimeInsert) (lastInsertId int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO animes(title, slug, abbreviation, kana, eng_name, description, thumb_url, state, series_id, count_episodes, copyright) "+
			"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		a.Title, a.Slug, a.Abbreviation, a.Kana, a.EngName, a.Description,
		a.ThumbUrl, a.State, a.SeriesId, a.CountEpisodes, a.Copyright,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawId, err := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	lastInsertId = int(rawId)
	return
}

func (repo *AdminAnimeRepository) Update(id int, a domain.TAnimeInsert) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"UPDATE animes SET title = ?, abbreviation = ?, kana = ?, eng_name = ?, description = ?, thumb_url = ?, state = ?, series_id = ?, count_episodes = ?, copyright = ? WHERE id = ?",
		a.Title, a.Abbreviation, a.Kana, a.EngName, a.Description, a.ThumbUrl, a.State, a.SeriesId, a.CountEpisodes, a.Copyright, id,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawAffected, err := exe.RowsAffected()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rowsAffected = int(rawAffected)
	return
}

func (repo *AdminAnimeRepository) Delete(id int) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"DELETE FROM animes WHERE id = ?", id,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawAffected, err := exe.RowsAffected()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rowsAffected = int(rawAffected)
	return
}

/************************
         platform
*************************/

func (repo *AdminPlatformRepository) ListAll() (platforms domain.TPlatforms, err error) {
	rows, err := repo.Query(
		"Select * from platforms",
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var p domain.TPlatform
		err = rows.Scan(
			&p.ID, &p.EngName, &p.PlatName,
			&p.BaseUrl, &p.Image, &p.IsValid,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
			return
		}
		platforms = append(platforms, p)
	}
	return
}

func (repo *AdminPlatformRepository) FindById(id int) (p domain.TPlatform, err error) {
	rows, err := repo.Query(
		"SELECT * FROM platforms WHERE id = ?", id,
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rows.Next()
	err = rows.Scan(
		&p.ID, &p.EngName, &p.PlatName, &p.BaseUrl,
		&p.Image, &p.IsValid, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	return
}

func (repo *AdminPlatformRepository) Insert(p domain.TPlatform) (lastInserted int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO platforms(eng_name, plat_name, base_url, image, is_valid) VALUES(?, ?, ?, ?, ?)",
		p.EngName, p.PlatName, p.BaseUrl, p.Image, p.IsValid,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawId, err := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	lastInserted = int(rawId)
	return
}

func (repo *AdminPlatformRepository) Update(p domain.TPlatform, id int) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"UPDATE platforms SET eng_name = ?, plat_name = ?, base_url = ?, image = ?, is_valid = ? WHERE id = ?",
		p.EngName, p.PlatName, p.BaseUrl, p.Image, p.IsValid, id,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawId, err := exe.RowsAffected()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rowsAffected = int(rawId)
	return
}

func (repo *AdminPlatformRepository) Delete(id int) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"DELETE FROM platforms WHERE id = ?", id,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawAffected, err := exe.RowsAffected()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rowsAffected = int(rawAffected)
	return
}

// relation

func (repo *AdminPlatformRepository) FilterByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	rows, err := repo.Query(
		"Select relation_anime_platform.*, platforms.plat_name FROM relation_anime_platform "+
			"LEFT JOIN platforms ON relation_anime_platform.platform_id = platforms.id "+
			"WHERE anime_id = ?", animeId,
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var p domain.TRelationPlatform
		err = rows.Scan(
			&p.PlatformId, &p.AnimeId, &p.LinkUrl,
			&p.CreatedAt, &p.UpdatedAt, &p.PlatName,
		)
		if err != nil {
			tools.ErrorLog(err)
			return
		}
		platforms = append(platforms, p)
	}
	return
}

func (repo *AdminPlatformRepository) InsertRelation(p domain.TRelationPlatformInput) (lastInserted int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO relation_anime_platform(platform_id, anime_id, link_url) VALUES(?, ?, ?)",
		p.PlatformId, p.AnimeId, p.LinkUrl,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawId, err := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	lastInserted = int(rawId)
	return
}

func (repo *AdminPlatformRepository) DeleteRelation(animeId int, platformId int) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"DELETE FROM relation_anime_platform WHERE anime_id = ? AND platform_id = ?",
		animeId, platformId,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawAffected, err := exe.RowsAffected()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rowsAffected = int(rawAffected)
	return
}

/************************
         season
*************************/

func (repo *AdminSeasonRepository) ListAll() (seasons []domain.TSeason, err error) {
	rows, err := repo.Query(
		"Select * from seasons",
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var s domain.TSeason
		err = rows.Scan(
			&s.ID, &s.Year, &s.Season, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		seasons = append(seasons, s)
	}
	return
}

func (repo *AdminSeasonRepository) FindById(id int) (s domain.TSeason, err error) {
	rows, err := repo.Query(
		"SELECT * FROM seasons WHERE id = ?", id,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(
		s.ID, s.Year, s.Season, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	return
}

func (repo *AdminSeasonRepository) Insert(s domain.TSeasonInput) (lastInserted int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO seasons(year, season) VALUES(?, ?)", s.Year, s.Season,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawId, _ := exe.LastInsertId()
	lastInserted = int(rawId)
	return
}

func (repo *AdminSeasonRepository) InsertRelation(r domain.TSeasonRelationInput) (lastInserted int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO relation_anime_season(season_id, anime_id) VALUES(?, ?)",
		r.SeasonId, r.AnimeId,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawId, _ := exe.LastInsertId()
	lastInserted = int(rawId)
	return
}

/************************
         series
*************************/

func (repo *AdminSeriesRepository) ListAll() (series []domain.TSeries, err error) {
	rows, err := repo.Query("Select * from series")
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var s domain.TSeries
		err = rows.Scan(
			&s.ID, &s.EngName, &s.SeriesName, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			tools.ErrorLog(err)
		}
		series = append(series, s)
	}
	return
}

func (repo *AdminSeriesRepository) FindById(id int) (s domain.TSeries, err error) {
	rows, err := repo.Query("SELECT * FROM series WHERE id = ?", id)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	defer rows.Close()
	rows.Next()
	err = rows.Scan(
		&s.ID, &s.EngName, &s.SeriesName, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		err = domain.ErrNotFound
		tools.ErrorLog(err)
		return
	}
	return
}

func (repo *AdminSeriesRepository) Insert(s domain.TSeriesInput) (lastInserted int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO series(eng_name, series_name) VALUES(?, ?)", s.EngName, s.SeriesName,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawId, _ := exe.LastInsertId()
	lastInserted = int(rawId)
	return
}

func (repo *AdminSeriesRepository) Update(s domain.TSeriesInput, id int) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"UPDATE series SET eng_name = ?, series_name = ? WHERE id = ?",
		s.EngName, s.SeriesName, id,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawId, _ := exe.LastInsertId()
	rowsAffected = int(rawId)
	return
}

func (repo *AdminSeriesRepository) Delete(id int) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"DELETE FROM series WHERE id = ?", id,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawId, _ := exe.LastInsertId()
	rowsAffected = int(rawId)
	return
}
