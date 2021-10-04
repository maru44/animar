package database

import (
	"animar/v1/internal/pkg/domain"

	"github.com/maru44/perr"
)

type AudienceRepository struct {
	SqlHandler
}

func (repo *AudienceRepository) Counts(animeId int) (audiences []domain.TAudienceCount, err error) {
	rows, err := repo.Query(
		"Select state, count(state) from audiences WHERE anime_id = ? GROUP BY state", animeId,
	)
	if err != nil {
		return audiences, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var a domain.TAudienceCount
		err := rows.Scan(
			&a.State, &a.Count,
		)
		if err != nil {
			return audiences, perr.Wrap(err, perr.NotFound)
		}
		audiences = append(audiences, a)
	}
	return
}

func (repo *AudienceRepository) FilterByUser(userId string) (audiences []domain.TAudienceJoinAnime, err error) {
	rows, err := repo.Query(
		"SELECT audiences.*, animes.title, animes.slug, animes.description, animes.state "+
			"FROM audiences LEFT JOIN animes ON audiences.anime_id = animes.id WHERE user_id = ?", userId,
	)
	if err != nil {
		return audiences, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var a domain.TAudienceJoinAnime
		err := rows.Scan(
			&a.ID, &a.State, &a.AnimeId, &a.UserId, &a.CreatedAt, &a.UpdatedAt,
			&a.Title, &a.Slug, &a.Content, &a.AState,
		)
		if err != nil {
			return audiences, perr.Wrap(err, perr.NotFound)
		}
		audiences = append(audiences, a)
	}
	return
}

func (repo *AudienceRepository) Insert(a domain.TAudienceInput, userId string) (lastInserted int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO audiences(state, anime_id, user_id) VALUES(?, ?, ?)",
		a.State, a.AnimeId, userId,
	)
	if err != nil {
		return lastInserted, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	rawId, err := exe.LastInsertId()
	if err != nil {
		return lastInserted, perr.Wrap(err, perr.BadRequest)
	}
	lastInserted = int(rawId)
	return
}

func (repo *AudienceRepository) Upsert(a domain.TAudienceInput, userId string) (rowsAffected int, err error) {
	_, err = repo.FindByAnimeAndUser(a.AnimeId, userId)
	if err == nil {
		exe, err := repo.Execute(
			"UPDATE audiences SET state = ? WHERE user_id = ? AND anime_id = ?",
			a.State, userId, a.AnimeId,
		)
		if err != nil {
			return rowsAffected, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
		}
		rawId, err := exe.RowsAffected()
		if err != nil {
			return rowsAffected, perr.Wrap(err, perr.BadRequest)
		}
		rowsAffected = int(rawId)
	} else {
		rowsAffected, err = repo.Insert(a, userId)
		if err != nil {
			return rowsAffected, perr.Wrap(err, perr.BadRequest)
		}
	}
	return
}

func (repo *AudienceRepository) Delete(animeId int, userId string) (rowsAffected int, err error) {
	exe, err := repo.Execute(
		"DELETE FROM audiences WHERE anime_id = ? AND user_id = ?",
		animeId, userId,
	)
	if err != nil {
		return rowsAffected, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	rawId, err := exe.RowsAffected()
	if err != nil {
		return rowsAffected, perr.Wrap(err, perr.BadRequest)
	}
	rowsAffected = int(rawId)
	return
}

func (repo *AudienceRepository) FindByAnimeAndUser(animeId int, userId string) (a domain.TAudience, err error) {
	rows, err := repo.Query(
		"Select * from audiences WHERE user_id = ? AND anime_id = ?",
		userId, animeId,
	)
	if err != nil {
		return a, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(
		&a.ID, &a.State, &a.AnimeId, &a.UserId, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return a, perr.Wrap(err, perr.NotFound)
	}
	return
}
