package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"net/http"
	"strconv"
)

type SeasonController struct {
	interactor domain.SeasonInteractor
}

func NewSeasonController(sqlHandler database.SqlHandler) *SeasonController {
	return &SeasonController{
		interactor: usecase.NewSeasonInteractor(
			&database.SeasonRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (controller *SeasonController) SeasonByAnimeIdView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	seasons, err := controller.interactor.RelationSeasonByAnime(id)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": seasons})
	return nil
}
