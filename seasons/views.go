package seasons

import (
	"animar/v1/tools/api"
	"animar/v1/tools/fire"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type TSeasonRelationInput struct {
	SeasonId int `json:"season_id"`
	AnimeId  int `json:"anime_id"`
}

func SeasonView(w http.ResponseWriter, r *http.Request) error {
	var userId string

	query := r.URL.Query()
	id := query.Get("id")

	if userId == "" {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("Forbidden")
	} else {
		if id != "" {
			i, _ := strconv.Atoi(id)
			s := detailSeason(i)
			if s.ID == 0 {
				w.WriteHeader(http.StatusNotFound)
				return errors.New("Not Found")
			}
			api.JsonResponse(w, map[string]interface{}{"data": s})
		} else {
			ss := ListSeasonDomain()
			api.JsonResponse(w, map[string]interface{}{"data": ss})
		}
	}
	return nil
}

func InsertSeasonView(w http.ResponseWriter, r *http.Request) error {
	var p TSeasonInput
	json.NewDecoder(r.Body).Decode(&p)
	insertedId := insertSeason(
		p.Year, p.Season,
	)
	api.JsonResponse(w, map[string]interface{}{"data": insertedId})
	return nil
}

/************************************
             relation
************************************/

func SeasonByAnimeIdView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	seasons := SeasonByAnimeIdDomain(id)
	api.JsonResponse(w, map[string]interface{}{"data": seasons})
	return nil
}

func InsertRelationSeasonView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetAdminIdFromCookie(r)

	if userId == "" {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("Forbidden")
	} else {
		var s TSeasonRelationInput
		json.NewDecoder(r.Body).Decode(&s)
		value := insertRelation(
			s.SeasonId, s.AnimeId,
		)
		api.JsonResponse(w, map[string]interface{}{"data": value})
	}
	return nil
}
