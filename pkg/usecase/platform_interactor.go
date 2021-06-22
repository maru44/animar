package usecase

import "animar/v1/pkg/domain"

type PlatformInteractor struct {
	repository PlatformRepository
}

func NewPlatformInteractor(platform PlatformRepository) domain.PlatformInteractor {
	return &PlatformInteractor{
		repository: platform,
	}
}

/************************
        repository
************************/

type PlatformRepository interface {
	FilterByAnime(int) (domain.TRelationPlatforms, error)
}

/**********************
   interactor methods
***********************/

func (interactor *PlatformInteractor) RelationPlatformByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	platforms, err = interactor.repository.FilterByAnime(animeId)
	return
}
