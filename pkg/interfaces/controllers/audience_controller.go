package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/interfaces/httphandle"
	"animar/v1/pkg/mvc/auth"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type AudienceController struct {
	interactor domain.AudienceInteractor
	//BaseController // うまくいかず
	base domain.BaseInteractor // うまくいかず
}

func NewAudienceController(sqlHandler database.SqlHandler) *AudienceController {
	return &AudienceController{
		interactor: usecase.NewAudienceInteractor(
			&database.AudienceRepository{
				SqlHandler: sqlHandler,
			},
		),
		// 以下機能せず
		base: usecase.NewBaseInteractor(
			&httphandle.BaseRepository{
				Firebase: infrastructure.NewFireBaseClient(),
			},
		),
	}
}

func (controller *AudienceController) AnimeAudienceCountsView(w http.ResponseWriter, r *http.Request) (ret error) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	audiences, err := controller.interactor.AnimeAudienceCounts(animeId)
	ret = response(w, err, map[string]interface{}{"data": audiences})
	return ret
}

func (controller *AudienceController) AudienceWithAnimeByUserView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := r.URL.Query().Get("user")
	audiences, err := controller.interactor.AudienceWithAnimeByUser(userId)
	ret = response(w, err, map[string]interface{}{"data": audiences})
	return ret
}

func (controller *AudienceController) UpsertAudienceView(w http.ResponseWriter, r *http.Request) (ret error) {
	/*  userId 取得  */
	idToken, _ := r.Cookie("idToken")
	claims := auth.VerifyFirebase(context.Background(), idToken.Value)
	userId := claims["user_id"].(string)

	if userId == "" {
		ret = response(w, domain.ErrUnauthorized, nil)
		return ret
	}
	var p domain.TAudienceInput
	json.NewDecoder(r.Body).Decode(&p)
	_, err := controller.interactor.UpsertAudience(p, userId)
	if err != nil {
		tools.ErrorLog(err)
	}
	ret = response(w, err, map[string]interface{}{"data": p.State})
	return ret
}

func (controller *AudienceController) DeleteAudienceView(w http.ResponseWriter, r *http.Request) (ret error) {
	/*  userId 取得  */
	idToken, _ := r.Cookie("idToken")
	claims := auth.VerifyFirebase(context.Background(), idToken.Value)
	userId := claims["user_id"].(string)

	if userId == "" {
		ret = response(w, domain.ErrUnauthorized, nil)
		return ret
	}

	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	_, err := controller.interactor.DeleteAudience(animeId, userId)
	ret = response(w, err, nil)
	return ret
}

func (controller *AudienceController) AudienceByAnimeAndUserView(w http.ResponseWriter, r *http.Request) (ret error) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	/*  userId 取得  */
	// userId, _ := controller.getUserIdFromCookie(r)
	//userId, _ := GetUserId(r)
	//userId, _ := controller.getUserIdFromCookie(r)
	idToken, _ := r.Cookie("idToken")
	claims := auth.VerifyFirebase(context.Background(), idToken.Value)
	userId := claims["user_id"].(string)

	if userId == "" {
		ret = response(w, domain.ErrUnauthorized, nil)
	} else {
		watch, _ := controller.interactor.AudienceByAnimeAndUser(animeId, userId)
		// 一旦 nil にしない。これはユーザーの視聴データが無いときにも対応するため
		ret = response(w, nil, map[string]interface{}{"data": watch})
	}
	return ret
}

// utility

func (controller *AudienceController) getUserIdFromCookie(r *http.Request) (userId string, err error) {
	idToken, err := r.Cookie("idToken")
	if err != nil {
		return
	} else if idToken.Value == "" {
		return
	}
	userId, err = controller.base.UserId(idToken.Value)
	return
}
