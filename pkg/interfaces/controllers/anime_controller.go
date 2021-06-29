package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/mvc/auth"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"context"
	"net/http"
	"strconv"
)

type AnimeController struct {
	interactor domain.AnimeInteractor
	BaseController
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

func (controller *AnimeController) AnimeListView(w http.ResponseWriter, r *http.Request) error {
	animes, _ := controller.interactor.AnimesAll()
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}

func (controller *AnimeController) AnimeView(w http.ResponseWriter, r *http.Request) (ret error) {
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
		ret = response(w, err, map[string]interface{}{"data": a})
	case slug != "":
		a, err := controller.interactor.AnimeDetailBySlug(slug)

		/*  userId 取得  */
		idToken, _ := r.Cookie("idToken")
		claims := auth.VerifyFirebase(context.Background(), idToken.Value)
		userId := claims["user_id"].(string)

		revs, _ := controller.interactor.ReviewFilterByAnime(a.GetId(), userId)
		ret = response(w, err, map[string]interface{}{"anime": a, "reviews": revs})
	case year != "":
		animes, err := controller.interactor.AnimesBySeason(year, season)
		ret = response(w, err, map[string]interface{}{"data": animes})
	case keyword != "":
		animes, err := controller.interactor.AnimesSearch(keyword)
		ret = response(w, err, map[string]interface{}{"data": animes})
	default:
		animes, err := controller.interactor.AnimesOnAir()
		ret = response(w, err, map[string]interface{}{"data": animes})
	}
	return ret
}

func (controller *AnimeController) SearchAnimeMinView(w http.ResponseWriter, r *http.Request) (ret error) {
	query := r.URL.Query()
	title := query.Get("t")

	animes, err := controller.interactor.AnimeSearchMinimum(title)
	if err != nil {
		tools.ErrorLog(err)
	}
	ret = response(w, err, map[string]interface{}{"data": animes})
	return ret
}

func (controller *AnimeController) AnimeMinimumsView(w http.ResponseWriter, r *http.Request) (ret error) {
	animes, err := controller.interactor.AnimeMinimums()
	if err != nil {
		tools.ErrorLog(err)
	}
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	ret = response(w, err, map[string]interface{}{"data": animes})
	return ret
}
