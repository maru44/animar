package database

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/tools/tools"
	"fmt"
)

type ReviewRepository struct {
	SqlHandler
}

func (repo *ReviewRepository) FindAll() (reviewIds []int, err error) {
	rows, err := repo.Query(
		"SELECT id FROM reviews",
	)
	defer rows.Close()

	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			domain.ErrorLog(err, "")
			return
		}
		reviewIds = append(reviewIds, id)
	}
	return
}

func (repo *ReviewRepository) FindById(id int) (r domain.ReviewWithAnimeSlug, err error) {
	rows, err := repo.Query(
		"SELECT r.*, a.slug, a.title FROM reviews AS r "+
			"LEFT JOIN animes AS a ON r.anime_id = a.id "+
			"WHERE r.id = ?", id,
	)
	defer rows.Close()

	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	rows.Next()
	err = rows.Scan(
		&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId, &r.CreatedAt, &r.UpdatedAt,
		&r.AnimeSlug, &r.AnimeTitle,
	)
	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	return
}

func (repo *ReviewRepository) FindByAnimeAndUser(animeId int, userId string) (r domain.TReview, err error) {
	rows, err := repo.Query(
		"SELECT * FROM reviews WHERE anime_id = ? AND user_id = ?", animeId, userId,
	)
	defer rows.Close()

	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	rows.Next()
	err = rows.Scan(
		&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId, &r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		domain.ErrorLog(err, "")
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
		domain.ErrorLog(err, "")
		return
	}
	for rows.Next() {
		var r domain.TReview
		err = rows.Scan(&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			domain.ErrorLog(err, "")
			return
		}
		reviews = append(reviews, r)
	}
	if err != nil {
		domain.ErrorLog(err, "")
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
		domain.ErrorLog(err, "")
		return
	}
	for rows.Next() {
		var r domain.TReviewJoinAnime
		err = rows.Scan(
			&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId,
			&r.CreatedAt, &r.UpdatedAt, &r.Title, &r.Slug, &r.AnimeContent, &r.AState,
		)
		if err != nil {
			domain.ErrorLog(err, "")
			return
		}
		reviews = append(reviews, r)
	}
	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	return
}

func (repo *ReviewRepository) InsertContent(r domain.TReviewInput, userId string) (content string, err error) {
	exe, err := repo.Execute(
		"INSERT INTO reviews(anime_id, content, user_id) VALUES(?, ?, ?)",
		r.AnimeId, r.Content, userId,
	)
	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	_, err = exe.LastInsertId()
	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	content = r.Content
	return
}

func (repo *ReviewRepository) UpsertContent(r domain.TReviewInput, userId string) (content string, err error) {
	review, err := repo.FindByAnimeAndUser(r.AnimeId, userId)
	if err == nil && review.GetId() != 0 {
		_, err = repo.Execute(
			"UPDATE reviews SET content = ? WHERE id = ?", r.Content, review.GetId(),
		)
		content = r.Content
		if err != nil {
			domain.ErrorLog(err, "")
			return
		}
		return
	}
	return repo.InsertContent(r, userId)
}

func (repo *ReviewRepository) InsertRating(r domain.TReviewInput, userId string) (rating int, err error) {
	exe, err := repo.Execute(
		"INSERT INTO reviews(anime_id, rating, user_id) VALUES(?, ?, ?)",
		r.AnimeId, r.Rating, userId,
	)
	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	_, err = exe.LastInsertId()
	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	rating = r.Rating
	return
}

func (repo *ReviewRepository) UpsertRating(r domain.TReviewInput, userId string) (rating int, err error) {
	review, err := repo.FindByAnimeAndUser(r.AnimeId, userId)
	if err == nil {
		_, err = repo.Execute(
			"UPDATE reviews SET rating = ? WHERE id = ?", tools.NewNullInt(r.Rating), review.GetId(),
		)
		if err != nil {
			domain.ErrorLog(err, "")
			return
		}
		rating = r.Rating
		return
	}
	return repo.InsertRating(r, userId)
}

func (repo *ReviewRepository) GetRatingAverage(animeId int) (rating string, err error) {
	rows, err := repo.Query("SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE anime_id = ?", animeId)
	defer rows.Close()
	if err != nil {
		domain.ErrorLog(err, "")
		return
	}
	var avg float32
	rows.Next()
	rows.Scan(&avg)
	return fmt.Sprintf("%.1f", avg), err
}
