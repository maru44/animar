package watch

import (
	"animar/v1/tools"
	"encoding/json"
	"net/http"
	"strconv"
)

type TWatchCountJsonResponse struct {
	Status int           `json:"Status"`
	Data   []TWatchCount `json:"Data"`
}

type TUserWatchJoinResponse struct {
	Status int               `json:"Status"`
	Data   []TWatchJoinAnime `json:"Data"`
}

type TWatchInput struct {
	AnimeId int    `json:"AnimeId"`
	Watch   int    `json:"Watch,string"` // form
	UserId  string `json:"UserId"`
}

func (result TWatchCountJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
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

func (result TUserWatchJoinResponse) ResponseWrite(w http.ResponseWriter) bool {
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

// by anime
func AnimeWatchCountView(w http.ResponseWriter, r *http.Request) error {
	result := TWatchCountJsonResponse{Status: 200}
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	var watches []TWatchCount
	watches = AnimeWatchCountDomain(animeId)

	result.Data = watches
	result.ResponseWrite(w)
	return nil
}

// by userId
// ?user=userId
func UserWatchStatusView(w http.ResponseWriter, r *http.Request) error {
	result := TUserWatchJoinResponse{Status: 200}
	userId := r.URL.Query().Get("user")

	var watches []TWatchJoinAnime
	watches = OnesWatchStatusDomain(userId)

	result.Data = watches
	result.ResponseWrite(w)
	return nil
}

// anime by ?anime=
// user by cookie
func WatchAnimeStateOfUserView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TIntJsonReponse{Status: 200}

	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	userId := tools.GetIdFromCookie(r)
	if userId == "" {
		result.Status = 4001
	} else {
		watch := WatchDetail(userId, animeId)
		result.Num = watch
	}

	result.ResponseWrite(w)
	return nil
}

// watch post view
func WatchPostView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetIdFromCookie(r)
	if userId == "" {
		result.Status = 4001
		result.ResponseWrite(w)
		return nil
	}

	var posted TWatchInput
	json.NewDecoder(r.Body).Decode(&posted)
	// watch := InsertWatch(posted.AnimeId, posted.Watch, userId)
	watch := UpsertWatch(posted.AnimeId, posted.Watch, userId)
	result.Num = watch

	result.ResponseWrite(w)
	return nil
}

// watch delete view
// ?anime=
func WatchDeleteView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TVoidJsonResponse{Status: 200}
	userId := tools.GetIdFromCookie(r)
	if userId == "" {
		result.Status = 5000
		result.ResponseWrite(w)
		return nil
	}

	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	exe := DeleteWatch(animeId, userId)
	if !exe {
		result.Status = 4000
	}

	result.ResponseWrite(w)
	return nil
}
