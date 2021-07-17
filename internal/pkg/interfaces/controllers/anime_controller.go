package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/tools/tools"
	"animar/v1/internal/pkg/usecase"
	"net/http"
	"strconv"
)

type AnimeController struct {
	interactor domain.AnimeInteractor
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
	}
}

func (controller *AnimeController) AnimeListView(w http.ResponseWriter, r *http.Request) {
	animes, err := controller.interactor.AnimesAll()
	response(w, err, map[string]interface{}{"data": animes})
	return
}

func (controller *AnimeController) AnimeView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	slug := query.Get("slug")
	year := query.Get("year")
	season := query.Get("season")
	keyword := query.Get("keyword")
	rawSeries := query.Get("series")

	switch {
	case strId != "":
		id, _ := strconv.Atoi(strId)
		a, err := controller.interactor.AnimeDetail(id)
		response(w, err, map[string]interface{}{"data": a})
	case slug != "":
		a, err := controller.interactor.AnimeDetailBySlug(slug)

		userId := r.Context().Value(USER_ID).(string)
		revs, _ := controller.interactor.ReviewFilterByAnime(a.GetId(), userId)
		response(w, err, map[string]interface{}{"anime": a, "reviews": revs})
	case year != "":
		animes, err := controller.interactor.AnimesBySeason(year, season)
		response(w, err, map[string]interface{}{"data": animes})
	case keyword != "":
		animes, err := controller.interactor.AnimesSearch(keyword)
		response(w, err, map[string]interface{}{"data": animes})
	case rawSeries != "":
		seriesId, _ := strconv.Atoi(rawSeries)
		animes, err := controller.interactor.AnimesBySeries(seriesId)
		response(w, err, map[string]interface{}{"data": animes})
	default:
		animes, err := controller.interactor.AnimesOnAir()
		response(w, err, map[string]interface{}{"data": animes})
	}
	return
}

func (controller *AnimeController) SearchAnimeMinimumView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	title := query.Get("t")

	animes, err := controller.interactor.AnimeSearchMinimum(title)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": animes})
	return
}

func (controller *AnimeController) AnimeMinimumsView(w http.ResponseWriter, r *http.Request) {
	animes, err := controller.interactor.AnimeMinimums()
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": animes})
	return
}
