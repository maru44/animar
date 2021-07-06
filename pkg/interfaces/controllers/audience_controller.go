package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"encoding/json"
	"net/http"
	"strconv"
)

type AudienceController struct {
	interactor domain.AudienceInteractor
}

func NewAudienceController(sqlHandler database.SqlHandler) *AudienceController {
	return &AudienceController{
		interactor: usecase.NewAudienceInteractor(
			&database.AudienceRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (controller *AudienceController) AnimeAudienceCountsView(w http.ResponseWriter, r *http.Request) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	audiences, err := controller.interactor.AnimeAudienceCounts(animeId)
	response(w, err, map[string]interface{}{"data": audiences})
	return
}

func (controller *AudienceController) AudienceWithAnimeByUserView(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user")
	audiences, err := controller.interactor.AudienceWithAnimeByUser(userId)
	response(w, err, map[string]interface{}{"data": audiences})
	return
}

func (controller *AudienceController) UpsertAudienceView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	var p domain.TAudienceInput
	json.NewDecoder(r.Body).Decode(&p)
	_, err := controller.interactor.UpsertAudience(p, userId)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": p.State})
	return
}

func (controller *AudienceController) DeleteAudienceView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	_, err := controller.interactor.DeleteAudience(animeId, userId)
	response(w, err, nil)
	return
}

func (controller *AudienceController) AudienceByAnimeAndUserView(w http.ResponseWriter, r *http.Request) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	userId := r.Context().Value(USER_ID).(string)
	watch, _ := controller.interactor.AudienceByAnimeAndUser(animeId, userId)
	// 一旦 nil にしない。これはユーザーの視聴データが無いときにも対応するため
	response(w, nil, map[string]interface{}{"data": watch})
	return
}
