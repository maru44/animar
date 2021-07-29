package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"net/http"
	"strconv"
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
	id, _ := strconv.Atoi(strId)

	platforms, err := controller.interactor.RelationPlatformByAnime(id)
	if err != nil {
		lg := domain.NewErrorLog(err.Error(), "")
		lg.Logging()
	}
	response(w, err, map[string]interface{}{"data": platforms})
	return
}
