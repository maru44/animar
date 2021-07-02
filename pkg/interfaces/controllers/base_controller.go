package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/fires"
	"animar/v1/pkg/usecase"
	"net/http"
)

type BaseController struct {
	interactor domain.BaseInteractor
}

func NewBaseController() *BaseController {
	return &BaseController{
		interactor: usecase.NewBaseInteractor(
			&fires.AuthRepository{
				Firebase: infrastructure.NewFireBaseClient(),
			},
		),
	}
}

func (controller *BaseController) getClaimsFromCookie(r *http.Request) (claims map[string]interface{}, err error) {
	idToken, err := r.Cookie("idToken")
	claims, err = controller.interactor.Claims(idToken.Value)
	return
}

func (controller *BaseController) getUserIdFromCookie(r *http.Request) (userId string, err error) {
	idToken, err := r.Cookie("idToken")
	if err != nil {
		return
	} else if idToken.Value == "" {
		return
	}
	userId, err = controller.interactor.UserId(idToken.Value)
	return
}

func (controller *BaseController) GetUserIdFromToken(idToken string) (userId string, err error) {
	claims, err := controller.interactor.Claims(idToken)
	if err != nil {
		return
	}
	userId = claims["user_id"].(string)
	return
}
