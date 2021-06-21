package usecase

import "animar/v1/pkg/domain"

type ReviewInteractor struct {
	ReviewRepository ReviewRepository
}

func (interactor *ReviewInteractor) GetOnesReviewByAnime(animeId int, userId string) (animes domain.TReview, err error) {
	animes, err = interactor.ReviewRepository.FindByAnimeAndUser(animeId, userId)
	return
}
func (interactor *ReviewInteractor) GetAnimeReviews(animeId int, userId string) (reviews domain.TReviews, err error) {
	reviews, err = interactor.ReviewRepository.FilterByAnime(animeId, userId)
	return
}
func (interactor *ReviewInteractor) GetOnesReviews(userId string) (reviews domain.TReviewJoinAnimes, err error) {
	reviews, err = interactor.ReviewRepository.FilterByUser(userId)
	return
}

func (interactor *ReviewInteractor) PostReviewContent(review domain.TReviewInput) (content string, err error) {
	content, err = interactor.ReviewRepository.InsertContent(review)
	return
}
func (interactor *ReviewInteractor) UpsertReviewContent(review domain.TReviewInput) (content string, err error) {
	content, err = interactor.ReviewRepository.UpsertContent(review)
	return
}

func (interactor *ReviewInteractor) PostReviewRating(review domain.TReviewInput) (rating int, err error) {
	rating, err = interactor.ReviewRepository.InsertRating(review)
	return
}
func (interactor *ReviewInteractor) UpsertReviewRating(review domain.TReviewInput) (rating int, err error) {
	rating, err = interactor.ReviewRepository.UpsertRating(review)
	return
}
func (interactor *ReviewInteractor) GetRatingAverage(animeId int) (rating string, err error) {
	rating, err = interactor.ReviewRepository.GetRatingAverage(animeId)
	return
}
