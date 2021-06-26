package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/apis"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type AudienceController struct {
	interactor domain.AudienceInteractor
	api        apis.ApiResponse
}

func NewAudienceController(sqlHandler database.SqlHandler) *AudienceController {
	return &AudienceController{
		interactor: usecase.NewAudienceInteractor(
			&database.AudienceRepository{
				SqlHandler: sqlHandler,
			},
		),
		api: infrastructure.NewApiResponse(),
	}
}

func (controller *AudienceController) AnimeAudienceCountsView(w http.ResponseWriter, r *http.Request) (ret error) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	audiences, err := controller.interactor.AnimeAudienceCounts(animeId)
	fmt.Println(audiences, err)
	ret = controller.api.Response(w, err, map[string]interface{}{"data": audiences})
	return ret
}

func (controller *AudienceController) AudienceWithAnimeByUserView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := r.URL.Query().Get("user")
	audiences, err := controller.interactor.AudienceWithAnimeByUser(userId)
	ret = controller.api.Response(w, err, map[string]interface{}{"data": audiences})
	return ret
}

func (controller *AudienceController) UpsertAudienceView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		ret = controller.api.Response(w, domain.ErrUnauthorized, nil)
		return ret
	}
	var p domain.TAudienceInput
	json.NewDecoder(r.Body).Decode(&p)
	_, err := controller.interactor.UpsertAudience(p, userId)
	if err != nil {
		tools.ErrorLog(err)
	}
	ret = controller.api.Response(w, err, map[string]interface{}{"data": p.State})
	return ret
}

func (controller *AudienceController) DeleteAudienceView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		ret = controller.api.Response(w, domain.ErrUnauthorized, nil)
		return ret
	}

	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	_, err := controller.interactor.DeleteAudience(animeId, userId)
	ret = controller.api.Response(w, err, nil)
	return ret
}

func (controller *AudienceController) AudienceByAnimeAndUserView(w http.ResponseWriter, r *http.Request) (ret error) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		ret = controller.api.Response(w, domain.ErrUnauthorized, nil)
	} else {
		watch, err := controller.interactor.AudienceByAnimeAndUser(animeId, userId)
		ret = controller.api.Response(w, err, map[string]interface{}{"data": watch})
	}
	return ret
}
