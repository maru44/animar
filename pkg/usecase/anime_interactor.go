package usecase

import "animar/v1/pkg/domain"

type AnimeInteractor struct {
	AnimeRepository AnimeRepository
}

func (interactor *AnimeInteractor) AnimesAll() (animes domain.TAnimes, err error) {
	animes, err = interactor.AnimeRepository.ListAll()
	return
}

func (interactor *AnimeInteractor) AnimesOnAir() (animes domain.TAnimes, err error) {
	animes, err = interactor.AnimeRepository.ListOnAirAll()
	return
}

func (interactor *AnimeInteractor) AnimesSearch(title string) (animes domain.TAnimes, err error) {
	animes, err = interactor.AnimeRepository.ListSearch(title)
	return
}

func (interactor *AnimeInteractor) AnimesBySeason(year string, season string) (animes domain.TAnimes, err error) {
	animes, err = interactor.AnimeRepository.ListBySeason(year, season)
	return
}

func (interactor *AnimeInteractor) AnimeSearchMinimum(title string) (animes domain.TAnimeMinimums, err error) {
	animes, err = interactor.AnimeRepository.ListMinimumSearch(title)
	return
}

// detail

func (interactor *AnimeInteractor) AnimeDetail(id int) (anime domain.TAnime, err error) {
	anime, err = interactor.AnimeRepository.FindById(id)
	return
}

func (interactor *AnimeInteractor) AnimeDetailBySlug(slug string) (anime domain.TAnimeWithSeries, err error) {
	anime, err = interactor.AnimeRepository.FindBySlug(slug)
	return
}
