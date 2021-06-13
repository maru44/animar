package seasons

import (
	"animar/v1/tools"
	"encoding/json"
	"net/http"
	"strconv"
)

type TSeasonRelationInput struct {
	SeasonId int `json:"season_id"`
	AnimeId  int `json:"anime_id"`
}

func SeasonView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	var userId string

	switch r.Method {
	case "GET":
		userId = tools.GetAdminIdFromCookie(r)
	case "POST":
		var posted tools.TUserIdCookieInput
		json.NewDecoder(r.Body).Decode(&posted)
		userId = tools.GetAdminIdFromIdToken(posted.Token)
	default:
		userId = ""
	}

	query := r.URL.Query()
	id := query.Get("id")

	if userId == "" {
		result.Status = 4003
	} else {
		var ss []TSeason
		if id != "" {
			i, _ := strconv.Atoi(id)
			s := detailSeason(i)
			if s.ID == 0 {
				result.Status = 404
			}
			ss = append(ss, s)
		} else {
			ss = ListSeasonDomain()
		}
		result.Data = ss
	}

	result.ResponseWrite(w)
	return nil
}

func InsertSeasonView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	var p TSeasonInput
	json.NewDecoder(r.Body).Decode(&p)
	insertedId := insertSeason(
		p.Year, p.Season,
	)
	result.Data = insertedId

	result.ResponseWrite(w)
	return nil
}

/************************************
             relation
************************************/

func SeasonByAnimeIdView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	seasons := SeasonByAnimeIdDomain(id)
	result.Data = seasons
	result.ResponseWrite(w)
	return nil
}

func InsertRelationSeasonView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	result, is_valid := result.LimitMethod([]string{"POST"}, r)
	if !is_valid {
		result.ResponseWrite(w)
		return nil
	}

	userId := tools.GetAdminIdFromCookie(r)

	if userId == "" {
		result.Status = 4003
	} else {
		var s TSeasonRelationInput
		json.NewDecoder(r.Body).Decode(&s)
		value := insertRelation(
			s.SeasonId, s.AnimeId,
		)
		result.Data = value
	}
	result.ResponseWrite(w)
	return nil
}
