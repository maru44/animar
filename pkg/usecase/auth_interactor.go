package usecase

import (
	"animar/v1/pkg/domain"
	"fmt"
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
	GetUserInfo(string) (domain.TUserInfo, error)
	GetClaims(string) (map[string]interface{}, error)
	IsAdmin(string) bool
	GetAdminId(string) (string, error)
	SendVerifyEmail(string) error
	Update(string, domain.TProfileForm) (domain.TUserInfo, error)
	GetUserId(string) (string, error)
}

/**********************
   interactor methods
***********************/

func (interactor *AuthInteractor) UserInfo(userId string) (user domain.TUserInfo, err error) {
	user, err = interactor.repository.GetUserInfo(userId)
	return
}

func (interactor *AuthInteractor) IsAdmin(userId string) bool {
	return interactor.repository.IsAdmin(userId)
}

func (interactor *AuthInteractor) AdminId(idToken string) (string, error) {
	fmt.Print("interactor")
	return interactor.repository.GetAdminId(idToken)
}

func (interactor *AuthInteractor) SendVerify(email string) error {
	return interactor.repository.SendVerifyEmail(email)
}

func (interactor *AuthInteractor) UpdateProfile(userId string, params domain.TProfileForm) (domain.TUserInfo, error) {
	return interactor.repository.Update(userId, params)
}

func (interactor *AuthInteractor) Claims(idToken string) (claims map[string]interface{}, err error) {
	return interactor.repository.GetClaims(idToken)
}
