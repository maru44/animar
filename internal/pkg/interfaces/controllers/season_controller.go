package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"net/http"
	"strconv"

	"github.com/maru44/perr"
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
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
	}

	seasons, err := controller.interactor.RelationSeasonByAnime(id)
	response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": seasons})
	return
}
