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

func AnimeView(w http.ResponseWriter, r *http.Request) error {
	// helper.SetDefaultResponseHeader(w)
	result := TAnimesJsonResponse{Status: 200}

	query := r.URL.Query()
	strId := query.Get("id")

	var animes []TAnime
	if strId != "" {
		id, err := strconv.Atoi(strId)
		if err != nil {
			panic(err.Error())
		}
		ani := DetailAnime(id)
		animes = append(animes, ani)
	} else {
		animes = ListAnimeDomain()
	}

	result.Data = animes
	result.ResponseWrite(w)

	return nil
}
