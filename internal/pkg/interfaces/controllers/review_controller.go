package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/maru44/perr"
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

func (controller *ReviewController) GetReviewView(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	if strId != "" {
		id, err := strconv.Atoi(strId)
		if err != nil {
			response(w, r, perr.Wrap(err, perr.BadRequest), nil)
			return
		}

		rev, err := controller.interactor.GetReviewById(id)
		response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": rev})
	} else {
		revs, err := controller.interactor.GetAllReviewIds()
		response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": revs})
	}
	return
}

func (controller *ReviewController) GetAnimeReviewsView(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user")
	animeIdStr := r.URL.Query().Get("anime")
	animeId, err := strconv.Atoi(animeIdStr)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	revs, err := controller.interactor.GetAnimeReviews(animeId, userId)
	response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": revs})
	return
}

func (controller *ReviewController) GetAnimeReviewOfUserView(w http.ResponseWriter, r *http.Request) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, err := strconv.Atoi(animeIdStr)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	userId := r.Context().Value(USER_ID).(string)
	rev, _ := controller.interactor.GetOnesReviewByAnime(animeId, userId)
	// 一旦 nil にしない。これはユーザーの視聴データが無いときにも対応するため
	response(w, r, nil, map[string]interface{}{"data": rev})
	return
}

func (controller *ReviewController) UpsertReviewContentView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	var posted domain.TReviewInput
	err := json.NewDecoder(r.Body).Decode(&posted)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	value, err := controller.interactor.UpsertReviewContent(posted, userId)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": value})
	return
}

func (controller *ReviewController) UpsertReviewRatingView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	var posted domain.TReviewInput
	err := json.NewDecoder(r.Body).Decode(&posted)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	value, err := controller.interactor.UpsertReviewRating(posted, userId)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": value})
	return
}

func (controller *ReviewController) AnimeRatingAvgView(w http.ResponseWriter, r *http.Request) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, err := strconv.Atoi(animeIdStr)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	avg, err := controller.interactor.GetRatingAverage(animeId)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": avg})
	return
}

func (controller *ReviewController) GetOnesReviewsView(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user")
	revs, err := controller.interactor.GetOnesReviews(userId)
	response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": revs})
	return
}
