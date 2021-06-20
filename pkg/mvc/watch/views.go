package watch

import (
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/fire"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type TAudienceInput struct {
	AnimeId int    `json:"anime_id"`
	State   int    `json:"state,string"` // form
	UserId  string `json:"user_id"`
}

// by anime
func AnimeWatchCountView(w http.ResponseWriter, r *http.Request) error {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	var watches []TAudienceCount
	watches = AnimeWatchCountDomain(animeId)

	api.JsonResponse(w, map[string]interface{}{"data": watches})
	return nil
}

// by userId
// ?user=userId
func UserWatchStatusView(w http.ResponseWriter, r *http.Request) error {
	userId := r.URL.Query().Get("user")

	var watches []TAudienceJoinAnime
	watches = OnesWatchStatusDomain(userId)

	api.JsonResponse(w, map[string]interface{}{"data": watches})
	return nil
}

// anime by ?anime=
// user by cookie
func WatchAnimeStateOfUserView(w http.ResponseWriter, r *http.Request) error {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorize")
	} else {
		watch := WatchDetail(userId, animeId)
		api.JsonResponse(w, map[string]interface{}{"data": watch})
	}
	return nil
}

// watch post view
func WatchPostView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorize")
	}

	var p TAudienceInput
	json.NewDecoder(r.Body).Decode(&p)
	// watch := InsertWatch(posted.AnimeId, posted.Watch, userId)
	watch := UpsertWatch(p.AnimeId, p.State, userId)
	api.JsonResponse(w, map[string]interface{}{"data": watch})
	return nil
}

// watch delete view
// ?anime=
func WatchDeleteView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorize")
	}

	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	exe := DeleteWatch(animeId, userId)
	if !exe {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Bad Request")
	}

	api.JsonResponse(w, map[string]interface{}{})
	return nil
}
