package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
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
			&database.CompanyRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (controller *AnimeController) AnimeListView(w http.ResponseWriter, r *http.Request) {
	animes, err := controller.interactor.AnimesAll()
	response(w, r, err, map[string]interface{}{"data": animes})
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
	company := query.Get("company")

	switch {
	case strId != "":
		id, err := strconv.Atoi(strId)
		if err != nil {
			response(w, r, domain.NewWrapError(err, domain.DataNotFoundError), nil)
			return
		}
		a, err := controller.interactor.AnimeDetail(id)
		if err != nil {
			response(w, r, err, nil)
			return
		}
		response(w, r, err, map[string]interface{}{"data": a})
	case slug != "":
		a, err := controller.interactor.AnimeDetailBySlug(slug)
		if err != nil {
			response(w, r, err, nil)
			return
		}

		userId := r.Context().Value(USER_ID).(string)
		revs, err := controller.interactor.ReviewFilterByAnime(a.GetId(), userId)
		response(w, r, err, map[string]interface{}{"anime": a, "reviews": revs})
	case year != "":
		animes, err := controller.interactor.AnimesBySeason(year, season)
		response(w, r, err, map[string]interface{}{"data": animes})
	case keyword != "":
		animes, err := controller.interactor.AnimesSearch(keyword)
		response(w, r, err, map[string]interface{}{"data": animes})
	case rawSeries != "":
		seriesId, err := strconv.Atoi(rawSeries)
		if err != nil {
			response(w, r, err, nil)
			return
		}
		animes, err := controller.interactor.AnimesBySeries(seriesId)
		response(w, r, err, map[string]interface{}{"data": animes})
	case company != "":
		comp, err := controller.interactor.DetailCompanyByEng(company)
		if err != nil {
			response(w, r, err, nil)
			return
		}
		animes, err := controller.interactor.AnimesByCompany(company)
		response(w, r, err, map[string]interface{}{"data": animes, "company": comp})
	default:
		animes, err := controller.interactor.AnimesOnAir()
		response(w, r, err, map[string]interface{}{"data": animes})
	}
	return
}

func (controller *AnimeController) SearchAnimeMinimumView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	title := query.Get("t")

	animes, err := controller.interactor.AnimeSearchMinimum(title)
	if err != nil {
		domain.ErrorWarn(err)
	}
	response(w, r, err, map[string]interface{}{"data": animes})
	return
}

func (controller *AnimeController) AnimeMinimumsView(w http.ResponseWriter, r *http.Request) {
	animes, err := controller.interactor.AnimeMinimums()
	if err != nil {
		domain.ErrorWarn(err)
	}
	response(w, r, err, map[string]interface{}{"data": animes})
	return
}
