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

func (controller *PlatformController) RelationPlatformByAnimeView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	platforms, err := controller.interactor.RelationPlatformByAnime(id)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": platforms})
	return nil
}
