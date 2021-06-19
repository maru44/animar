package controllers

import (
	"animar/v1/tools/api"
	"animar/v1/usecase"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type AnimeController struct {
	// interactorがnull
	Interactor usecase.AnimeInteractor
}

func NewAnimeController() *AnimeController {
	return &AnimeController{}
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
		a, err := controller.Interactor.AnimeDetail(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": a})
	case slug != "":
		a, err := controller.Interactor.AnimeDetailBySlug(slug)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		// @TODO review
		api.JsonResponse(w, map[string]interface{}{"data": a})
	case year != "":
		animes, _ := controller.Interactor.AnimesBySeason(year, season)
		api.JsonResponse(w, map[string]interface{}{"data": animes})
	case keyword != "":
		animes, _ := controller.Interactor.AnimesSearch(keyword)
		api.JsonResponse(w, map[string]interface{}{"data": animes})
	default:
		//animes, _ := controller.Interactor.AnimesOnAir()
		fmt.Println(controller)
		animes, _ := controller.Interactor.AnimesAll()
		api.JsonResponse(w, map[string]interface{}{"data": animes})
	}
	return nil
}

func (controller *AnimeController) ListAnimeMinimumView(w http.ResponseWriter, r *http.Request) error {
	animes, _ := controller.Interactor.AnimesMinimumAll()
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}
