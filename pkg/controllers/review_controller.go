package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/apis"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/usecase"
	"encoding/json"
	"net/http"
	"strconv"
)

type ReviewController struct {
	interactor domain.ReviewInteractor
	api        apis.ApiResponse
}

func NewReviewController(sqlHandler database.SqlHandler) *ReviewController {
	return &ReviewController{
		interactor: usecase.NewReviewInteractor(
			&database.ReviewRepository{
				SqlHandler: sqlHandler,
			},
		),
		api: infrastructure.NewApiResponse(),
	}
}

func (controller *ReviewController) GetAnimeReviewsView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := r.URL.Query().Get("user")
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	revs, err := controller.interactor.GetAnimeReviews(animeId, userId)

	ret = controller.api.Response(w, err, map[string]interface{}{"data": revs})
	return ret
}

func (controller *ReviewController) GetAnimeReviewOfUserView(w http.ResponseWriter, r *http.Request) (ret error) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		ret = controller.api.Response(w, domain.ErrUnauthorized, nil)
	} else {
		rev, err := controller.interactor.GetOnesReviewByAnime(animeId, userId)
		ret = controller.api.Response(w, err, map[string]interface{}{"data": rev})
	}
	return ret
}

func (controller *ReviewController) UpsertReviewContentView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		ret = controller.api.Response(w, domain.ErrUnauthorized, nil)
	} else {
		var posted domain.TReviewInput
		json.NewDecoder(r.Body).Decode(&posted)
		value, err := controller.interactor.UpsertReviewContent(posted, userId)
		ret = controller.api.Response(w, err, map[string]interface{}{"data": value})
	}
	return ret
}

func (controller *ReviewController) UpsertReviewRatingView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		ret = controller.api.Response(w, domain.ErrUnauthorized, nil)
	} else {
		var posted domain.TReviewInput
		json.NewDecoder(r.Body).Decode(&posted)
		value, err := controller.interactor.UpsertReviewRating(posted, userId)
		ret = controller.api.Response(w, err, map[string]interface{}{"data": value})
	}
	return ret
}

func (controller *ReviewController) AnimeRatingAvgView(w http.ResponseWriter, r *http.Request) (ret error) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	avg, err := controller.interactor.GetRatingAverage(animeId)
	ret = controller.api.Response(w, err, map[string]interface{}{"data": avg})
	return ret
}

func (controller *ReviewController) GetOnesReviewsView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := r.URL.Query().Get("user")
	revs, err := controller.interactor.GetOnesReviews(userId)
	ret = controller.api.Response(w, err, map[string]interface{}{"data": revs})
	return ret
}
