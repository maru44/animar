package usecase

import (
	"animar/v1/pkg/domain"
)

type BaseInteractor struct {
	repository BaseRepository
}

func NewBaseInteractor(b BaseRepository) domain.BaseInteractor {
	return &BaseInteractor{
		repository: b,
	}
}

/************************
        repository
************************/

type BaseRepository interface {
	GetUserId(string) (string, error)
	GetClaims(string) (map[string]interface{}, error)
}

/**********************
   interactor methods
***********************/

func (interactor *BaseInteractor) UserId(idToken string) (string, error) {
	return interactor.repository.GetUserId(idToken)
}

func (interactor *BaseInteractor) Claims(idToken string) (claims map[string]interface{}, err error) {
	claims, err = interactor.repository.GetClaims(idToken)
	return
}
