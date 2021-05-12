package review

import (
	"animar/v1/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

/*
func GetOnesReviews(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	userId :=
}
*/

type TReviewsJsonResponse struct {
	Status int       `json:"Status"`
	Data   []TReview `json:"Data"`
}

type TReviewInput struct {
	AnimeId int    `json:"AnimeId,string"`
	Content string `json:"Content"`
	Star    int    `json:"Star,string"`
	UserId  string `json:"UserId"`
}

func (result TReviewsJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	helper.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func GetYourReviews(w http.ResponseWriter, r *http.Request) error {
	result := TReviewsJsonResponse{Status: 200}

	userId := helper.GetIdFromCookie(r)
	if userId == "" {
		result.Status = 5000
		result.ResponseWrite(w)
		return nil
	}

	var revs []TReview
	revs = OnesReviewsDomain(userId)

	result.Data = revs
	result.ResponseWrite(w)
	return nil
}

func GetAnimeReviews(w http.ResponseWriter, r *http.Request) error {
	result := TReviewsJsonResponse{Status: 200}
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	var revs []TReview
	revs = AnimeReviewsDomain(animeId)

	result.Data = revs
	result.ResponseWrite(w)
	return nil
}

func ReviewPostJsonView(w http.ResponseWriter, r *http.Request) error {
	result := helper.TIntJsonReponse{Status: 200}

	userId := helper.GetIdFromCookie(r)
	if userId == "" {
		result.Status = 5000
		result.ResponseWrite(w)
		return nil
	}

	var posted TReviewInput
	json.NewDecoder(r.Body).Decode(&posted)
	insertedId := InsertReview(posted.AnimeId, posted.Content, posted.Star, userId)

	result.Num = insertedId

	fmt.Print(userId)
	result.ResponseWrite(w)
	return nil
}

func ReviewPostView(w http.ResponseWriter, r *http.Request) error {
	result := helper.TIntJsonReponse{Status: 200}

	userId := helper.GetIdFromCookie(r)
	if userId == "" {
		result.Status = 5000
		result.ResponseWrite(w)
		return nil
	}

	var posted TReviewInput
	json.NewDecoder(r.Body).Decode(&posted)
	insertedId := InsertReview(posted.AnimeId, posted.Content, posted.Star, userId)

	result.Num = insertedId

	result.ResponseWrite(w)
	return nil
}

func ReviewPostSample(w http.ResponseWriter, r *http.Request) error {
	result := helper.TIntJsonReponse{Status: 200}

	var posted TReviewInput
	json.NewDecoder(r.Body).Decode(&posted)

	// nullAnimeId = new(sql.NullInt32)
	// nullStar = new(sql.NullInt32)
	fmt.Print(posted)
	fmt.Print(reflect.TypeOf(posted.AnimeId))

	result.ResponseWrite(w)
	return nil
}

func ReviewTest(w http.ResponseWriter, r *http.Request) error {
	result := TReviewsJsonResponse{Status: 200}
	userId := helper.GetIdFromCookie(r)

	fmt.Print(userId)
	result.ResponseWrite(w)
	return nil
}
