package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"net/http"
	"strconv"

	"github.com/maru44/perr"
)

type PlatformController struct {
	interactor domain.PlatformInteractor
}

func NewPlatformController(sqlHandler database.SqlHandler) *PlatformController {
	return &PlatformController{
		interactor: usecase.NewPlatformInteractor(
			&database.PlatformRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (controller *PlatformController) RelationPlatformByAnimeView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
	}

	platforms, err := controller.interactor.RelationPlatformByAnime(id)
	response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": platforms})
	return
}
