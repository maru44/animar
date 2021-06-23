package database

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/tools/tools"
	"fmt"
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

func (repo *ReviewRepository) InsertContent(r domain.TReviewInput, userId string) (content string, err error) {
	exe, err := repo.Execute(
		"INSERT INTO reviews(anime_id, content, user_id) VALUES(?, ?, ?, ?)",
		r.AnimeId, r.Content, userId,
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

func (repo *ReviewRepository) UpsertContent(r domain.TReviewInput, userId string) (content string, err error) {
	review, err := repo.FindByAnimeAndUser(r.AnimeId, userId)
	if err == nil && review.GetId() != 0 {
		_, err = repo.Execute(
			"UPDATE reviews SET content = ? WHERE id = ?", r.Content, review.GetId(),
		)
		content = *r.Content
		if err != nil {
			tools.ErrorLog(err)
			return
		}
	}
	return repo.InsertContent(r, userId)
}

func (repo *ReviewRepository) InsertRating(r domain.TReviewInput, userId string) (rating int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO reviews(anime_id, rating, user_id) VALUES(?, ?, ?, ?)",
		r.AnimeId, r.Rating, userId,
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
	rating = *r.Rating
	return
}

func (repo *ReviewRepository) UpsertRating(r domain.TReviewInput, userId string) (rating int, err error) {
	review, err := repo.FindByAnimeAndUser(r.AnimeId, userId)
	if err == nil && review.GetId() != 0 {
		_, err = repo.Execute(
			"UPDATE reviews SET rating = ? WHERE id = ?", r.Rating, review.GetId(),
		)
		rating = *r.Rating
		if err != nil {
			tools.ErrorLog(err)
			return
		}
	}
	return repo.InsertRating(r, userId)
}

func (repo *ReviewRepository) GetRatingAverage(animeId int) (rating string, err error) {
	rows, err := repo.Query("SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE anime_id = ?", animeId)
	defer rows.Close()
	if err != nil {
		tools.ErrorLog(err)
		return
	}
	var avg float32
	rows.Next()
	rows.Scan(&avg)
	return fmt.Sprintf("%.1f", avg), err
}
