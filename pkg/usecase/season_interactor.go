package usecase

import "animar/v1/pkg/domain"

type SeasonInteractor struct {
	repository SeasonRepository
}

func NewSeasonInteractor(season SeasonRepository) domain.SeasonInteractor {
	return &SeasonInteractor{
		repository: season,
	}
}

/************************
        repository
************************/

type SeasonRepository interface {
	FilterByAnimeId(int) ([]domain.TSeasonRelation, error)
}

func (interactor *SeasonInteractor) RelationSeasonByAnime(animeId int) (seasons []domain.TSeasonRelation, err error) {
	seasons, err = interactor.repository.FilterByAnimeId(animeId)
	return
}
