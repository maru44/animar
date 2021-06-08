package seasons

import (
	"animar/v1/tools"
	"encoding/json"
	"net/http"
	"strconv"
)

type TSeasonRelationResponse struct {
	Status int               `json:"status"`
	Data   []TSeasonRelation `json:"data"`
}

type TSeasonResponse struct {
	Status int       `json:"status"`
	Data   []TSeason `json:"data"`
}

func (result TSeasonRelationResponse) ResponseWrite(w http.ResponseWriter) bool {
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

func (result TSeasonResponse) ResponseWrite(w http.ResponseWriter) bool {
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

func SeasonView(w http.ResponseWriter, r *http.Request) error {
	result := TSeasonResponse{Status: 200}
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
			s := DetailSeason(i)
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
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetAdminIdFromCookie(r)
	if userId == "" {
		result.Status = 4003
	} else {
		var p TSeasonInput
		json.NewDecoder(r.Body).Decode(&p)
		insertedId := InsertSeason(
			p.Year, p.Season,
		)
		result.Num = insertedId
	}
	result.ResponseWrite(w)
	return nil
}

/************************************
             relation
************************************/

func SeasonByAnimeIdView(w http.ResponseWriter, r *http.Request) error {
	result := TSeasonRelationResponse{Status: 200}

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	seasons := SeasonByAnimeIdDomain(id)
	result.Data = seasons
	result.ResponseWrite(w)
	return nil
}
