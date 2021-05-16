package anime

import (
	"animar/v1/helper"
	"encoding/json"
	"net/http"
	"strconv"
)

type TAnimesJsonResponse struct {
	Status int      `json:"Status"`
	Data   []TAnime `json:"Data"`
}

type TAnimesWithUserWatchResponse struct {
	Status int                   `json:"Status"`
	Data   []TAnimeWithUserWatch `json:"Data"`
}

type TAnimeInput struct {
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	OnAirState int    `json:"OnAirState"`
}

func (animeJson TAnimesJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(animeJson)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	helper.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func (animeWUWCJson TAnimesWithUserWatchResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(animeWUWCJson)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	helper.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

// list and detail
func AnimeView(w http.ResponseWriter, r *http.Request) error {
	result := TAnimesJsonResponse{Status: 200}

	query := r.URL.Query()
	strId := query.Get("id")
	slug := query.Get("slug")

	var animes []TAnime
	if strId != "" {
		id, _ := strconv.Atoi(strId)
		ani := DetailAnime(id)
		if ani.ID == 0 {
			result.Status = 404
		}
		animes = append(animes, ani)
	} else if slug != "" {
		ani := DetailAnimeBySlug(slug)
		if ani.ID == 0 {
			result.Status = 404
		}
		animes = append(animes, ani)
	} else {
		animes = ListAnimeDomain()
	}

	result.Data = animes
	result.ResponseWrite(w)

	return nil
}

// anime data + user's watch
func AnimeWithUserWatchView(w http.ResponseWriter, r *http.Request) error {
	result := TAnimesWithUserWatchResponse{Status: 200}

	query := r.URL.Query()
	strId := query.Get("id")
	slug := query.Get("slug")
	userId := helper.GetIdFromCookie(r)

	var animes []TAnimeWithUserWatch
	if strId != "" {
		/*
			id, _ := strconv.Atoi(strId)
			ani := DetailAnime(id)
			if ani.ID == 0 {
				result.Status = 404
			}
			animes = append(animes, ani)
		*/
	} else if slug != "" {
		ani := DetailAnimeBySlugWithUserWatch(slug, userId)
		if ani.ID == 0 {
			result.Status = 404
		}
		animes = append(animes, ani)
	} else {
		// animes = ListAnimeDomain()
	}

	result.Data = animes
	result.ResponseWrite(w)

	return nil
}

func AnimePostView(w http.ResponseWriter, r *http.Request) error {
	result := helper.TIntJsonReponse{Status: 200}

	var posted TAnimeInput
	json.NewDecoder(r.Body).Decode(&posted)
	insertedId := InsertAnime(posted.Title, posted.Content, posted.OnAirState)

	result.Num = insertedId
	result.ResponseWrite(w)

	return nil
}
