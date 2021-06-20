package controllers

import (
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/usecase"
	"animar/v1/tools/api"
	"net/http"
)

type AnimeController struct {
	Interactor usecase.AnimeInteractor
}

func NewAnimeController(sqlHandler database.SqlHandler) *AnimeController {
	return &AnimeController{
		Interactor: usecase.AnimeInteractor{
			AnimeRepository: &database.AnimeRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *AnimeController) AnimeListView(w http.ResponseWriter, r *http.Request) error {
	animes, _ := controller.Interactor.AnimesAll()
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}
