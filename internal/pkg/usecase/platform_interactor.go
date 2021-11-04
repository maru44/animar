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
	// slack notification
	RegisterTarget(in domain.NotifiedTargetInput) (int, error)
	UpdateTarget(in domain.NotifiedTargetInput) (int, error)
	DeleteTarget(userId string) (int, error)
	GetUsersChannel(userId string) (*string, error)
}

// batch

type PlatformBatchRepository interface {
	FilterNotificationTarget() ([]string, error)
	FilterTodaysBroadCast() ([]domain.NotificationBroadcast, error)
	MakeSlackMessage([]domain.NotificationBroadcast) string
	ChangeOnAirState(state string, slug string) (int, error)
}

/**********************
   interactor methods
***********************/

func (in *PlatformInteractor) RelationPlatformByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	platforms, err = in.repository.FilterByAnime(animeId)
	return
}

func (in *PlatformInteractor) RegisterNotifiedTarget(input domain.NotifiedTargetInput) (int, error) {
	return in.repository.RegisterTarget(input)
}

func (in *PlatformInteractor) UpdateNotifiedTarget(input domain.NotifiedTargetInput) (int, error) {
	return in.repository.UpdateTarget(input)
}

func (in *PlatformInteractor) DeleteNotifiedTarget(userId string) (int, error) {
	return in.repository.DeleteTarget(userId)
}

func (in *PlatformInteractor) UsersChannel(userId string) (*string, error) {
	return in.repository.GetUsersChannel(userId)
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

func (in *PlatformBatchInteractor) ChangeOnAirState(state string, slug string) (int, error) {
	return in.repo.ChangeOnAirState(state, slug)
}
