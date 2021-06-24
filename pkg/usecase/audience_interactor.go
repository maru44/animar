package usecase

import "animar/v1/pkg/domain"

type AudienceInteractor struct {
	repository AudienceRepository
}

func NewAudienceInteractor(a AudienceRepository) domain.AudienceInteractor {
	return &AudienceInteractor{
		repository: a,
	}
}

/************************
        repository
************************/

type AudienceRepository interface {
	Counts(int) ([]domain.TAudienceCount, error)
	FilterByUser(string) ([]domain.TAudienceJoinAnime, error)
	FindByAnimeAndUser(int, string) (domain.TAudience, error)
	Insert(domain.TAudienceInput, string) (int, error)
	Upsert(domain.TAudienceInput, string) (int, error)
	Delete(int, string) (int, error)
}

func (interactor *AudienceInteractor) AnimeAudienceCounts(animeId int) (audiences []domain.TAudienceCount, err error) {
	audiences, err = interactor.repository.Counts(animeId)
	return
}

func (interactor *AudienceInteractor) AudienceWithAnimeByUser(userId string) (audiences []domain.TAudienceJoinAnime, err error) {
	audiences, err = interactor.repository.FilterByUser(userId)
	return
}

func (interactor *AudienceInteractor) AudienceByAnimeAndUser(animeId int, userId string) (audience domain.TAudience, err error) {
	audience, err = interactor.repository.FindByAnimeAndUser(animeId, userId)
	return
}

func (interactor *AudienceInteractor) InsertAudience(a domain.TAudienceInput, userId string) (lastInserted int, err error) {
	lastInserted, err = interactor.repository.Insert(a, userId)
	return
}

func (interactor *AudienceInteractor) UpsertAudience(a domain.TAudienceInput, userId string) (rowsAffected int, err error) {
	rowsAffected, err = interactor.repository.Upsert(a, userId)
	return
}

func (interactor *AudienceInteractor) DeleteAudience(animeId int, userId string) (rowsAffected int, err error) {
	rowsAffected, err = interactor.repository.Delete(animeId, userId)
	return
}
