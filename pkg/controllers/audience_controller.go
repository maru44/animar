package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"encoding/json"
	"errors"
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

func (controller *AudienceController) AnimeAudienceCountsView(w http.ResponseWriter, r *http.Request) error {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	audiences, _ := controller.interactor.AnimeAudienceCounts(animeId)
	api.JsonResponse(w, map[string]interface{}{"data": audiences})
	return nil
}

func (controller *AudienceController) AudienceWithAnimeByUserView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorize")
	} else {
		audiences, _ := controller.interactor.AudienceWithAnimeByUser(userId)
		api.JsonResponse(w, map[string]interface{}{"data": audiences})
		return nil
	}
}

func (controller *AudienceController) UpsertAudienceView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorize")
	}
	var p domain.TAudienceInput
	json.NewDecoder(r.Body).Decode(&p)
	_, err := controller.interactor.UpsertAudience(p, userId)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": p.State})
	return nil
}

func (controller *AudienceController) DeleteAudience(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorize")
	}

	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	_, _ = controller.interactor.DeleteAudience(animeId, userId)
	api.JsonResponse(w, map[string]interface{}{})
	return nil
}
