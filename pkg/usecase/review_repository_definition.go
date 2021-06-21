package usecase

import "animar/v1/pkg/domain"

type ReviewRepository interface {
	FindByAnimeAndUser(int, string) (domain.TReview, error)
	FilterByAnime(int, string) (domain.TReviews, error)
	FilterByUser(string) (domain.TReviewJoinAnimes, error)
	// content
	InsertContent(domain.TReviewInput) (string, error)
	UpsertContent(domain.TReviewInput) (string, error)
	// Rating
	InsertRating(domain.TReviewInput) (int, error)
	UpsertRating(domain.TReviewInput) (int, error)
	GetRatingAverage(int) (string, error)
}
