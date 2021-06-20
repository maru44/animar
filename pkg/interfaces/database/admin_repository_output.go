package database

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/tools/tools"
)

type AdminRepository struct {
	SqlHandler
}

func (repo *AdminRepository) ListAll() (animes domain.TAnimes, err error) {
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

func (repo *AdminRepository) FindById(id int) (a domain.TAnimeAdmin, err error) {
	row, err := repo.Query(
		"SELECT * FROM animes WHERE id = ?", id,
	)
	defer row.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	row.Next()
	err = row.Scan(
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

func (repo *AdminRepository) Insert(a domain.TAnimeInsert) (lastInsertId int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO animes(title, slug, abbreviation, kana, eng_name, description, thumb_url, state, series_id, count_episodes, copyright) " +
			"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
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

func (repo *AdminRepository) Update(id int, a domain.TAnimeInsert) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"UPDATE animes SET title = ?, abbreviation = ?, kana = ?, eng_name = ?, description = ?, thumb_url = ?, state = ?, series_id = ?, count_episodes = ?, copyright = ? WHERE id = ?",
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawAffected, err := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rowsAffected = int(rawAffected)
	return
}

func (repo *AdminRepository) Delete(id int) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"DELETE FROM animes WHERE id = ?",
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rawAffected, err := exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rowsAffected = int(rawAffected)
	return
}
