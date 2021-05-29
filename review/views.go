package review

import (
	"animar/v1/tools"
	"encoding/json"
	"net/http"
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

type TReviewsWithUserInfoResponse struct {
	Status int               `json:"Status"`
	Data   []TReviewJoinUser `json:"Data"`
}

type TReviewInput struct {
	AnimeId int    `json:"AnimeId"`
	Content string `json:"Content"`
	Star    int    `json:"Star,string"` // text/plainのpostに対応
	UserId  string `json:"UserId"`
}

type TReviewJoinAnimeResponse struct {
	Status int                `json:"Status"`
	Data   []TReviewJoinAnime `json:"Data"`
}

func (result TReviewsJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	tools.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func (result TReviewJoinAnimeResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	tools.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func (result TReviewsWithUserInfoResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	tools.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

// cookie
func GetYourReviews(w http.ResponseWriter, r *http.Request) error {
	result := TReviewsJsonResponse{Status: 200}

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
	result := TReviewJoinAnimeResponse{Status: 200}

	userId := r.URL.Query().Get("user")
	var revs []TReviewJoinAnime
	revs = OnesReviewsJoinAnimeDomain(userId)

	result.Data = revs
	result.ResponseWrite(w)
	return nil
}

func GetAnimeReviews(w http.ResponseWriter, r *http.Request) error {
	result := TReviewsJsonResponse{Status: 200}
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
	result := TReviewsWithUserInfoResponse{Status: 200}
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
	result := TReviewsJsonResponse{Status: 200}

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
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetIdFromCookie(r)

	if userId == "" {
		result.Status = 4001
	} else {
		var posted TReviewInput
		json.NewDecoder(r.Body).Decode(&posted)
		value := UpsertReviewStar(posted.AnimeId, posted.Star, userId)
		result.Num = value
	}
	result.ResponseWrite(w)
	return nil
}

//upsert content
func UpsertReviewContentView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TStringJsonResponse{Status: 200}
	userId := tools.GetIdFromCookie(r)

	if userId == "" {
		result.Status = 4001
	} else {
		var posted TReviewInput
		json.NewDecoder(r.Body).Decode(&posted)
		upsertedString := UpsertReviewContent(posted.AnimeId, posted.Content, userId)
		result.String = upsertedString
	}
	result.ResponseWrite(w)
	return nil
}

// anime star avarage view
func AnimeStarAvgView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TStringJsonResponse{Status: 200}
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	avg := AnimeStarAvg(animeId)
	result.String = avg
	result.ResponseWrite(w)
	return nil
}
