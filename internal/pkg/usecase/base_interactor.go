package usecase

import (
	"animar/v1/internal/pkg/domain"
)

type BaseInteractor struct {
	authRepo AuthRepository
}

func NewBaseInteractor(a AuthRepository) domain.BaseInteractor {
	return &BaseInteractor{
		authRepo: a,
	}
}

/************************
        repository
************************/

// type BaseRepository interface {
// 	GetUserId(string) (string, error)
// 	GetClaims(string) (map[string]interface{}, error)
// }

/**********************
   interactor methods
***********************/

func (interactor *BaseInteractor) UserId(idToken string) (string, error) {
	return interactor.authRepo.GetUserId(idToken)
}

func (interactor *BaseInteractor) Claims(idToken string) (claims map[string]interface{}, err error) {
	claims, err = interactor.authRepo.GetClaims(idToken)
	return
}

func (interactor *BaseInteractor) AdminId(idToken string) (string, error) {
	return interactor.authRepo.GetAdminId(idToken)
}
