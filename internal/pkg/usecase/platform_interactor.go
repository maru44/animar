package usecase

import "animar/v1/internal/pkg/domain"

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
	// relation
	FilterTodaysBroadCast() ([]domain.NotificationBroadcast, error)
}

/**********************
   interactor methods
***********************/

func (in *PlatformInteractor) RelationPlatformByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	platforms, err = in.repository.FilterByAnime(animeId)
	return
}

func (in *PlatformInteractor) TargetNotificationBroadcast() ([]domain.NotificationBroadcast, error) {
	return in.repository.FilterTodaysBroadCast()
}
