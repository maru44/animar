package usecase

import "animar/v1/internal/pkg/domain"

type AnimeInteractor struct {
	animeRepository   AnimeRepository
	reviewRepository  ReviewAnimeRepository
	companyRepository CompanyRepository
}

func NewAnimeInteractor(anime AnimeRepository, review ReviewAnimeRepository, cop CompanyRepository) domain.AnimeInteractor {
	return &AnimeInteractor{
		animeRepository:   anime,
		reviewRepository:  review,
		companyRepository: cop,
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
	ListByCompany(string) (domain.TAnimes, error)
	// detail
	FindById(int) (domain.TAnime, error)
	FindBySlug(string) (domain.TAnimeWithSeries, error)
	// change on air state
}

type ReviewAnimeRepository interface {
	// review
	FilterByAnime(int, string) (domain.TReviews, error)
}

/**********************
   interactor methods
***********************/

func (in *AnimeInteractor) AnimesAll() (animes domain.TAnimes, err error) {
	animes, err = in.animeRepository.ListAll()
	return
}

func (in *AnimeInteractor) AnimesOnAir() (animes domain.TAnimes, err error) {
	animes, err = in.animeRepository.ListOnAirAll()
	return
}

func (in *AnimeInteractor) AnimeMinimums() (animes domain.TAnimeMinimums, err error) {
	animes, err = in.animeRepository.ListMinimum()
	return
}

func (in *AnimeInteractor) AnimesSearch(title string) (animes domain.TAnimes, err error) {
	animes, err = in.animeRepository.ListSearch(title)
	return
}

func (in *AnimeInteractor) AnimesBySeason(year string, season string) (animes domain.TAnimes, err error) {
	animes, err = in.animeRepository.ListBySeason(year, season)
	return
}

func (in *AnimeInteractor) AnimeSearchMinimum(title string) (animes domain.TAnimeMinimums, err error) {
	animes, err = in.animeRepository.ListMinimumSearch(title)
	return
}

func (in *AnimeInteractor) AnimesBySeries(id int) (animes []domain.TAnimeWithSeries, err error) {
	animes, err = in.animeRepository.ListBySeries(id)
	return
}

func (in *AnimeInteractor) AnimesByCompany(engName string) (domain.TAnimes, error) {
	return in.animeRepository.ListByCompany(engName)
}

// detail

func (in *AnimeInteractor) AnimeDetail(id int) (anime domain.TAnime, err error) {
	anime, err = in.animeRepository.FindById(id)
	return
}

func (in *AnimeInteractor) AnimeDetailBySlug(slug string) (anime domain.TAnimeWithSeries, err error) {
	anime, err = in.animeRepository.FindBySlug(slug)
	return
}

// review

func (in *AnimeInteractor) ReviewFilterByAnime(animeId int, userId string) (reviews domain.TReviews, err error) {
	reviews, err = in.reviewRepository.FilterByAnime(animeId, userId)
	return
}

// company

func (in *AnimeInteractor) DetailCompanyByEng(engName string) (domain.CompanyDetail, error) {
	return in.companyRepository.DetailByEng(engName)
}
