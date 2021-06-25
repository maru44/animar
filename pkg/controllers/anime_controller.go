package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/apis"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"errors"
	"net/http"
	"strconv"
)

type AnimeController struct {
	interactor domain.AnimeInteractor
	api        apis.ApiResponse
}

func NewAnimeController(sqlHandler database.SqlHandler) *AnimeController {
	return &AnimeController{
		interactor: usecase.NewAnimeInteractor(
			&database.AnimeRepository{
				SqlHandler: sqlHandler,
			},
			&database.ReviewRepository{
				SqlHandler: sqlHandler,
			},
		),
		api: infrastructure.NewApiResponse(),
	}
}

func (controller *AnimeController) AnimeListView(w http.ResponseWriter, r *http.Request) error {
	animes, _ := controller.interactor.AnimesAll()
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}

func (controller *AnimeController) AnimeView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	slug := query.Get("slug")
	year := query.Get("year")
	season := query.Get("season")
	keyword := query.Get("keyword")

	switch {
	case strId != "":
		id, _ := strconv.Atoi(strId)
		a, err := controller.interactor.AnimeDetail(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": a})
	case slug != "":
		a, err := controller.interactor.AnimeDetailBySlug(slug)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		userId := fire.GetIdFromCookie(r)
		revs, _ := controller.interactor.ReviewFilterByAnime(a.GetId(), userId)
		api.JsonResponse(w, map[string]interface{}{"anime": a, "reviews": revs})
	case year != "":
		animes, _ := controller.interactor.AnimesBySeason(year, season)
		api.JsonResponse(w, map[string]interface{}{"data": animes})
	case keyword != "":
		animes, _ := controller.interactor.AnimesSearch(keyword)
		api.JsonResponse(w, map[string]interface{}{"data": animes})
	default:
		animes, _ := controller.interactor.AnimesOnAir()
		api.JsonResponse(w, map[string]interface{}{"data": animes})
	}
	return nil
}

func (controller *AnimeController) SearchAnimeMinView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	title := query.Get("t")

	animes, err := controller.interactor.AnimeSearchMinimum(title)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}

func (controller *AnimeController) AnimeMinimumsView(w http.ResponseWriter, r *http.Request) error {
	animes, err := controller.interactor.AnimeMinimums()
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}
