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

type PlatformBatchInteractor struct {
	repo PlatformBatchRepository
}

func NewPlatformBatchInteractor(platform PlatformBatchRepository) domain.PlatformBatchInteractor {
	return &PlatformBatchInteractor{
		repo: platform,
	}
}

/************************
        repository
************************/

type PlatformRepository interface {
	FilterByAnime(int) (domain.TRelationPlatforms, error)
}

// batch

type PlatformBatchRepository interface {
	FilterNotificationTarget() ([]string, error)
	FilterTodaysBroadCast() ([]domain.NotificationBroadcast, error)
	MakeSlackMessage([]domain.NotificationBroadcast) string
}

/**********************
   interactor methods
***********************/

func (in *PlatformInteractor) RelationPlatformByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	platforms, err = in.repository.FilterByAnime(animeId)
	return
}

// batch

func (in *PlatformBatchInteractor) FilterNotificationTarget() ([]string, error) {
	return in.repo.FilterNotificationTarget()
}

func (in *PlatformBatchInteractor) TargetNotificationBroadcast() ([]domain.NotificationBroadcast, error) {
	return in.repo.FilterTodaysBroadCast()
}

func (in *PlatformBatchInteractor) MakeSlackMessage(nbs []domain.NotificationBroadcast) string {
	return in.repo.MakeSlackMessage(nbs)
}
