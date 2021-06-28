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
	IsAdmin(string) bool
	GetAdminId(context.Context, string) (string, error)
	SendVerifyEmail(context.Context, string) error
	Update(context.Context, string, domain.TProfileForm) (domain.TUserInfo, error)
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
	return interactor.repository.IsAdmin(userId)
}

func (interactor *AuthInteractor) AdminId(ctx context.Context, idToken string) (string, error) {
	return interactor.repository.GetAdminId(ctx, idToken)
}

func (interactor *AuthInteractor) SendVerify(ctx context.Context, email string) error {
	return interactor.repository.SendVerifyEmail(ctx, email)
}

func (interactor *AuthInteractor) UpdateProfile(ctx context.Context, userId string, params domain.TProfileForm) (domain.TUserInfo, error) {
	return interactor.repository.Update(ctx, userId, params)
}
