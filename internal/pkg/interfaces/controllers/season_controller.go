package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
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

func (controller *SeasonController) SeasonByAnimeIdView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	seasons, err := controller.interactor.RelationSeasonByAnime(id)
	if err != nil {
		lg := domain.NewErrorLog(err.Error(), "")
		lg.Logging()
	}
	response(w, err, map[string]interface{}{"data": seasons})
}
