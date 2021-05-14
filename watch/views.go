package watch

import (
	"animar/v1/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type TWatchCountJsonResponse struct {
	Status int           `json:"Status"`
	Data   []TWatchCount `json:"Data"`
}

type TUserWatchStatusResponse struct {
	Status int      `json:"Status"`
	Data   []TWatch `json:"Data"`
}

func (result TWatchCountJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
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

func (result TUserWatchStatusResponse) ResponseWrite(w http.ResponseWriter) bool {
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

func GetWatchCountOfAnime(w http.ResponseWriter, r *http.Request) {
	animeIdStr := r.URL.Query().Get("anime")
	animeId, _ := strconv.Atoi(animeIdStr)

	var watches []TWatchCount
	watches = AnimeWatchCountDomain(animeId)

	fmt.Print(watches)
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
func UserWatchStatusView(w http.ResponseWriter, r *http.Request) error {
	result := TUserWatchStatusResponse{Status: 200}
	userId := r.URL.Query().Get("user")

	var watches []TWatch
	watches = OnesWatchStatusDomain(userId)

	result.Data = watches
	result.ResponseWrite(w)
	return nil
}
