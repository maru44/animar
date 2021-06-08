package review

import (
	"animar/v1/tools"
	"encoding/json"
	"net/http"
	"strconv"
)

type TReviewInput struct {
	AnimeId int    `json:"anime_id"`
	Content string `json:"content,omitempty"`
	Star    int    `json:"rating,string,omitempty"` // text/plainのpostに対応
	UserId  string `json:"user_id"`
}

// cookie
func GetYourReviews(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	userId := tools.GetIdFromCookie(r)
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

func GetOnesReviewsView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	userId := r.URL.Query().Get("user")
	var revs []TReviewJoinAnime
	revs = OnesReviewsJoinAnimeDomain(userId)

	result.Data = revs
	result.ResponseWrite(w)
	return nil
}

func GetAnimeReviewsView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	userId := tools.GetIdFromCookie(r)

	var revs []TReview
	revs = AnimeReviewsDomain(animeId, userId)

	result.Data = revs
	result.ResponseWrite(w)
	return nil
}

func GetAnimeReviewsWithUserInfoView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)
	userId := tools.GetIdFromCookie(r)

	var revs []TReviewJoinUser
	revs = AnimeReviewsWithUserInfoDomain(animeId, userId)

	result.Data = revs
	result.ResponseWrite(w)
	return nil
}

// user's anime's review
func GetAnimeUserReviewView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	userId := tools.GetIdFromCookie(r)
	if userId == "" {
		result.Status = 4001
	} else {
		var revs []TReview
		rev := DetailReviewAnimeUser(animeId, userId)
		revs = append(revs, rev)
		result.Data = revs
	}
	result.ResponseWrite(w)
	return nil
}

// upsert star
func UpsertReviewStarView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	userId := tools.GetIdFromCookie(r)

	if userId == "" {
		result.Status = 4001
	} else {
		var posted TReviewInput
		json.NewDecoder(r.Body).Decode(&posted)
		value := UpsertReviewStar(posted.AnimeId, posted.Star, userId)
		result.Data = value
	}
	result.ResponseWrite(w)
	return nil
}

//upsert content
func UpsertReviewContentView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	userId := tools.GetIdFromCookie(r)

	if userId == "" {
		result.Status = 4001
	} else {
		var posted TReviewInput
		json.NewDecoder(r.Body).Decode(&posted)
		upsertedString := UpsertReviewContent(posted.AnimeId, posted.Content, userId)
		result.Data = upsertedString
	}
	result.ResponseWrite(w)
	return nil
}

// anime star avarage view
func AnimeStarAvgView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	avg := AnimeStarAvg(animeId)
	result.Data = avg
	result.ResponseWrite(w)
	return nil
}
