package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/mvc/auth"
	"animar/v1/pkg/usecase"
	"context"
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

func (controller *ReviewController) GetAnimeReviewsView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := r.URL.Query().Get("user")
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	revs, err := controller.interactor.GetAnimeReviews(animeId, userId)

	ret = response(w, err, map[string]interface{}{"data": revs})
	return ret
}

func (controller *ReviewController) GetAnimeReviewOfUserView(w http.ResponseWriter, r *http.Request) (ret error) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	/*  userId 取得  */
	idToken, _ := r.Cookie("idToken")
	claims := auth.VerifyFirebase(context.Background(), idToken.Value)
	userId := claims["user_id"].(string)

	if userId == "" {
		ret = response(w, domain.ErrUnauthorized, nil)
	} else {
		rev, _ := controller.interactor.GetOnesReviewByAnime(animeId, userId)
		// 一旦 nil にしない。これはユーザーの視聴データが無いときにも対応するため
		ret = response(w, nil, map[string]interface{}{"data": rev})
	}
	return ret
}

func (controller *ReviewController) UpsertReviewContentView(w http.ResponseWriter, r *http.Request) (ret error) {
	/*  userId 取得  */
	idToken, _ := r.Cookie("idToken")
	claims := auth.VerifyFirebase(context.Background(), idToken.Value)
	userId := claims["user_id"].(string)

	if userId == "" {
		ret = response(w, domain.ErrUnauthorized, nil)
	} else {
		var posted domain.TReviewInput
		json.NewDecoder(r.Body).Decode(&posted)
		value, err := controller.interactor.UpsertReviewContent(posted, userId)
		ret = response(w, err, map[string]interface{}{"data": value})
	}
	return ret
}

func (controller *ReviewController) UpsertReviewRatingView(w http.ResponseWriter, r *http.Request) (ret error) {
	/*  userId 取得  */
	idToken, _ := r.Cookie("idToken")
	claims := auth.VerifyFirebase(context.Background(), idToken.Value)
	userId := claims["user_id"].(string)

	if userId == "" {
		ret = response(w, domain.ErrUnauthorized, nil)
	} else {
		var posted domain.TReviewInput
		json.NewDecoder(r.Body).Decode(&posted)
		value, err := controller.interactor.UpsertReviewRating(posted, userId)
		ret = response(w, err, map[string]interface{}{"data": value})
	}
	return ret
}

func (controller *ReviewController) AnimeRatingAvgView(w http.ResponseWriter, r *http.Request) (ret error) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	avg, err := controller.interactor.GetRatingAverage(animeId)
	ret = response(w, err, map[string]interface{}{"data": avg})
	return ret
}

func (controller *ReviewController) GetOnesReviewsView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := r.URL.Query().Get("user")
	revs, err := controller.interactor.GetOnesReviews(userId)
	ret = response(w, err, map[string]interface{}{"data": revs})
	return ret
}
