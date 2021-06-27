package usecase

import (
	"animar/v1/pkg/domain"
	"context"
)

type AuthInteractor struct {
	repository AuthRepository
}

func NewAuthInteractor(auth AuthRepository) domain.AuthInteractor {
	return &AuthInteractor{
		repository: auth,
	}
}

/************************
        repository
************************/

type AuthRepository interface {
	GetUserInfo(context.Context, string) (domain.TUserInfo, error)
	GetClaims(context.Context, string) (map[string]interface{}, error)
	IsAmin(string) bool
	GetAdminId(context.Context, string) string
}

/**********************
   interactor methods
***********************/

func (interactor *AuthInteractor) UserInfo(ctx context.Context, userId string) (user domain.TUserInfo, err error) {
	user, err = interactor.repository.GetUserInfo(ctx, userId)
	return
}

func (interactor *AuthInteractor) Claims(ctx context.Context, idToken string) (claims map[string]interface{}, err error) {
	claims, err = interactor.repository.GetClaims(ctx, idToken)
	return
}

func (interactor *AuthInteractor) IsAdmin(userId string) bool {
	return interactor.repository.IsAmin(userId)
}

func (interactor *AuthInteractor) AdminId(ctx context.Context, idToken string) string {
	return interactor.repository.GetAdminId(ctx, idToken)
}
