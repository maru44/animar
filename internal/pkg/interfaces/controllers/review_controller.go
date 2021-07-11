package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"encoding/json"
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

func (controller *ReviewController) GetAnimeReviewsView(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user")
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	revs, err := controller.interactor.GetAnimeReviews(animeId, userId)
	response(w, err, map[string]interface{}{"data": revs})
	return
}

func (controller *ReviewController) GetAnimeReviewOfUserView(w http.ResponseWriter, r *http.Request) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	userId := r.Context().Value(USER_ID).(string)

	rev, _ := controller.interactor.GetOnesReviewByAnime(animeId, userId)
	// 一旦 nil にしない。これはユーザーの視聴データが無いときにも対応するため
	response(w, nil, map[string]interface{}{"data": rev})
	return
}

func (controller *ReviewController) UpsertReviewContentView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	var posted domain.TReviewInput
	json.NewDecoder(r.Body).Decode(&posted)
	value, err := controller.interactor.UpsertReviewContent(posted, userId)
	response(w, err, map[string]interface{}{"data": value})
	return
}

func (controller *ReviewController) UpsertReviewRatingView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	var posted domain.TReviewInput
	json.NewDecoder(r.Body).Decode(&posted)
	value, err := controller.interactor.UpsertReviewRating(posted, userId)
	response(w, err, map[string]interface{}{"data": value})
	return
}

func (controller *ReviewController) AnimeRatingAvgView(w http.ResponseWriter, r *http.Request) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	avg, err := controller.interactor.GetRatingAverage(animeId)
	response(w, err, map[string]interface{}{"data": avg})
	return
}

func (controller *ReviewController) GetOnesReviewsView(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user")
	revs, err := controller.interactor.GetOnesReviews(userId)
	response(w, err, map[string]interface{}{"data": revs})
	return
}
