package database

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/tools/tools"
)

type ReviewRepository struct {
	SqlHandler
}

func (repo *ReviewRepository) FindByAnimeAndUser(animeId int, userId string) (r domain.TReview, err error) {
	rows, err := repo.Query(
		"SELECT * FROM reviews WHERE anime_id = ? AND user_id = ?", animeId, userId,
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	rows.Next()
	err = rows.Scan(
		&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	return
}

func (repo *ReviewRepository) FilterByAnime(animeId int, userId string) (reviews domain.TReviews, err error) {
	rows, err := repo.Query(
		"Select * from reviews WHERE anime_id = ? AND (user_id != ? OR user_id IS NULL)",
		animeId, userId,
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var r domain.TReview
		err = rows.Scan(&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			tools.ErrorLog(err)
			return
		}
		reviews = append(reviews, r)
	}
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	return
}

func (repo *ReviewRepository) FilterByUser(userId string) (reviews domain.TReviewJoinAnimes, err error) {
	rows, err := repo.Query(
		"SELECT reviews.*, animes.title, animes.slug, animes.description, animes.state "+
			"FROM reviews LEFT JOIN animes ON reviews.anime_id = animes.id WHERE user_id = ?", userId,
	)
	defer rows.Close()

	if err != nil {
		tools.ErrorLog(err)
		return
	}
	for rows.Next() {
		var r domain.TReviewJoinAnime
		err = rows.Scan(
			&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId,
			&r.CreatedAt, &r.UpdatedAt, &r.Title, &r.Slug, &r.AnimeContent, &r.AState,
		)
		if err != nil {
			tools.ErrorLog(err)
			return
		}
		reviews = append(reviews, r)
	}
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	return
}

func (repo *ReviewRepository) InsertContent(r domain.TReviewInput) (content string, err error) {
	exe, err := repo.Execute(
		"INSERT INTO reviews(anime_id, content, rating, user_id) VALUES(?, ?, ?, ?)",
		r.AnimeId, r.Content, r.Rating, r.UserId,
	)
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	_, err = exe.LastInsertId()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	content = *r.Content
	return
}
