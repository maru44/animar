package usecase

import "animar/v1/internal/pkg/domain"

type AnimeInteractor struct {
	animeRepository  AnimeRepository
	reviewRepository ReviewAnimeRepository
}

func NewAnimeInteractor(anime AnimeRepository, review ReviewAnimeRepository) domain.AnimeInteractor {
	return &AnimeInteractor{
		animeRepository:  anime,
		reviewRepository: review,
	}
}

/************************
        repository
************************/

type AnimeRepository interface {
	ListAll() (domain.TAnimes, error)
	ListOnAirAll() (domain.TAnimes, error)
	ListMinimumSearch(string) (domain.TAnimeMinimums, error)
	ListMinimum() (domain.TAnimeMinimums, error)
	ListSearch(string) (domain.TAnimes, error)
	ListBySeason(string, string) (domain.TAnimes, error)
	ListBySeries(int) ([]domain.TAnimeWithSeries, error)
	// detail
	FindById(int) (domain.TAnime, error)
	FindBySlug(string) (domain.TAnimeWithSeries, error)
}

type ReviewAnimeRepository interface {
	// review
	FilterByAnime(int, string) (domain.TReviews, error)
}

/**********************
   interactor methods
***********************/

func (interactor *AnimeInteractor) AnimesAll() (animes domain.TAnimes, err error) {
	animes, err = interactor.animeRepository.ListAll()
	return
}

func (interactor *AnimeInteractor) AnimesOnAir() (animes domain.TAnimes, err error) {
	animes, err = interactor.animeRepository.ListOnAirAll()
	return
}

func (interactor *AnimeInteractor) AnimeMinimums() (animes domain.TAnimeMinimums, err error) {
	animes, err = interactor.animeRepository.ListMinimum()
	return
}

func (interactor *AnimeInteractor) AnimesSearch(title string) (animes domain.TAnimes, err error) {
	animes, err = interactor.animeRepository.ListSearch(title)
	return
}

func (interactor *AnimeInteractor) AnimesBySeason(year string, season string) (animes domain.TAnimes, err error) {
	animes, err = interactor.animeRepository.ListBySeason(year, season)
	return
}

func (interactor *AnimeInteractor) AnimeSearchMinimum(title string) (animes domain.TAnimeMinimums, err error) {
	animes, err = interactor.animeRepository.ListMinimumSearch(title)
	return
}

func (interactor *AnimeInteractor) AnimesBySeries(id int) (animes []domain.TAnimeWithSeries, err error) {
	animes, err = interactor.animeRepository.ListBySeries(id)
	return
}

// detail

func (interactor *AnimeInteractor) AnimeDetail(id int) (anime domain.TAnime, err error) {
	anime, err = interactor.animeRepository.FindById(id)
	return
}

func (interactor *AnimeInteractor) AnimeDetailBySlug(slug string) (anime domain.TAnimeWithSeries, err error) {
	anime, err = interactor.animeRepository.FindBySlug(slug)
	return
}

// review

func (interactor *AnimeInteractor) ReviewFilterByAnime(animeId int, userId string) (reviews domain.TReviews, err error) {
	reviews, err = interactor.reviewRepository.FilterByAnime(animeId, userId)
	return
}
