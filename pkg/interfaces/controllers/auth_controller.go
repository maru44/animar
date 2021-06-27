package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/fires"
	"animar/v1/pkg/usecase"
	"context"
	"net/http"
)

type AuthController struct {
	interactor domain.AuthInteractor
}

func NewAuthController(firebase fires.Firebase) *AuthController {
	return &AuthController{
		interactor: usecase.NewAuthInteractor(
			&fires.AuthRepository{
				Firebase: firebase,
			},
		),
	}
}

func (controller *AuthController) GetUserFromQueryView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := r.URL.Query().Get("uid")
	ctx := context.Background()

	user, err := controller.interactor.UserInfo(ctx, userId)
	ret = response(w, err, map[string]interface{}{"user": user})
	return
}

func (controller *AuthController) GetUserFromCookieView(w http.ResponseWriter, r *http.Request) (ret error) {
	ctx := context.Background()
	claims, err := controller.getClaimsFromCookie(r, ctx)
	userId := claims["user_id"].(string)
	user, err := controller.interactor.UserInfo(ctx, userId)
	ret = response(w, err, map[string]interface{}{"user": user, "is_verify": claims["email_verified"]})
	return
}

func (controller *AuthController) getClaimsFromCookie(r *http.Request, ctx context.Context) (claims map[string]interface{}, err error) {
	idToken, err := r.Cookie("idToken")
	claims, err = controller.interactor.Claims(ctx, idToken.Value)
	return
}
