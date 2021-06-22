package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type ReviewController struct {
	interactor domain.ReviewInteractor
}

func NewReviewController(sqlHandler database.SqlHandler) *ReviewController {
	return &ReviewController{
		interactor: usecase.NewReviewInteractor(
			&database.ReviewRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (controller *ReviewController) GetAnimeReviewsView(w http.ResponseWriter, r *http.Request) error {
	userId := r.URL.Query().Get("user")
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	revs, err := controller.interactor.GetAnimeReviews(animeId, userId)
	if err != nil {
		tools.ErrorLog(err)
		w.WriteHeader(http.StatusNotFound)
		return err
	}

	api.JsonResponse(w, map[string]interface{}{"data": revs})
	return nil
}

func (controller *ReviewController) GetAnimeReviewOfUserView(w http.ResponseWriter, r *http.Request) error {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized")
	} else {
		rev, err := controller.interactor.GetOnesReviewByAnime(animeId, userId)
		if err != nil {
			tools.ErrorLog(err)
			w.WriteHeader(http.StatusNotFound)
			return err
		}
		api.JsonResponse(w, map[string]interface{}{"data": rev})
		return nil
	}
}

func (controller *ReviewController) UpsertReviewContentView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized")
	} else {
		var posted domain.TReviewInput
		json.NewDecoder(r.Body).Decode(&posted)
		value, _ := controller.interactor.UpsertReviewContent(posted)
		api.JsonResponse(w, map[string]interface{}{"data": value})
		return nil
	}
}

func (controller *ReviewController) UpsertReviewRatingView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized")
	} else {
		var posted domain.TReviewInput
		value, _ := controller.interactor.UpsertReviewRating(posted)
		api.JsonResponse(w, map[string]interface{}{"data": value})
		return nil
	}
}

func (controller *ReviewController) AnimeRatingAvgView(w http.ResponseWriter, r *http.Request) error {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	avg, _ := controller.interactor.GetRatingAverage(animeId)
	api.JsonResponse(w, map[string]interface{}{"data": avg})
	return nil
}

func (controller *ReviewController) GetOnesReviewsView(w http.ResponseWriter, r *http.Request) error {
	userId := r.URL.Query().Get("user")
	revs, _ := controller.interactor.GetOnesReviews(userId)
	api.JsonResponse(w, map[string]interface{}{"data": revs})
	return nil
}
