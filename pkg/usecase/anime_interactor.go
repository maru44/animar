package usecase

import "animar/v1/pkg/domain"

type AnimeInteractor struct {
	AnimeRepository AnimeRepository
}

func (interactor *AnimeInteractor) AnimesAll() (animes domain.TAnimes, err error) {
	animes, err = interactor.AnimeRepository.ListAll()
	return
}
