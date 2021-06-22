package usecase

import "animar/v1/pkg/domain"

type ReviewInteractor struct {
	repository ReviewRepository
}

func NewReviewInteractor(review ReviewRepository) domain.ReviewInteractor {
	return &ReviewInteractor{
		repository: review,
	}
}

/************************
        repository
************************/

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

/**********************
   interactor methods
***********************/

func (interactor *ReviewInteractor) GetOnesReviewByAnime(animeId int, userId string) (animes domain.TReview, err error) {
	animes, err = interactor.repository.FindByAnimeAndUser(animeId, userId)
	return
}

func (interactor *ReviewInteractor) GetAnimeReviews(animeId int, userId string) (reviews domain.TReviews, err error) {
	reviews, err = interactor.repository.FilterByAnime(animeId, userId)
	return
}

func (interactor *ReviewInteractor) GetOnesReviews(userId string) (reviews domain.TReviewJoinAnimes, err error) {
	reviews, err = interactor.repository.FilterByUser(userId)
	return
}

func (interactor *ReviewInteractor) PostReviewContent(review domain.TReviewInput) (content string, err error) {
	content, err = interactor.repository.InsertContent(review)
	return
}
func (interactor *ReviewInteractor) UpsertReviewContent(review domain.TReviewInput) (content string, err error) {
	content, err = interactor.repository.UpsertContent(review)
	return
}

func (interactor *ReviewInteractor) PostReviewRating(review domain.TReviewInput) (rating int, err error) {
	rating, err = interactor.repository.InsertRating(review)
	return
}

func (interactor *ReviewInteractor) UpsertReviewRating(review domain.TReviewInput) (rating int, err error) {
	rating, err = interactor.repository.UpsertRating(review)
	return
}

func (interactor *ReviewInteractor) GetRatingAverage(animeId int) (rating string, err error) {
	rating, err = interactor.repository.GetRatingAverage(animeId)
	return
}
