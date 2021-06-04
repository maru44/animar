package series

import (
	"animar/v1/tools"
	"encoding/json"
	"net/http"
	"strconv"
)

type TSeriesResponse struct {
	Status int       `json:"Status"`
	Data   []TSeries `json:"Data"`
}

func (result TSeriesResponse) ResponseWrite(w http.ResponseWriter) bool {
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

func SeriesView(w http.ResponseWriter, r *http.Request) error {
	result := TSeriesResponse{Status: 200}
	var userId string

	switch r.Method {
	case "GET":
		userId = tools.GetAdminIdFromCookie(r)
	case "POST":
		var posted tools.TUserIdCookieInput
		json.NewDecoder(r.Body).Decode(&posted)
		userId = tools.GetAdminIdFromIdToken(posted.UserId)
	default:
		userId = ""
	}

	query := r.URL.Query()
	id := query.Get("id")

	if userId == "" {
		result.Status = 4003
	} else {
		var sers []TSeries
		if id != "" {
			i, _ := strconv.Atoi(id)
			ser := DetailSeries(i)
			if ser.ID == 0 {
				result.Status = 404
			}
			sers = append(sers, ser)
		} else {
			sers = ListSeriesDomain()
		}
		result.Data = sers
	}

	result.ResponseWrite(w)
	return nil
}

func InsertSeriesView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetAdminIdFromCookie(r)
	if userId == "" {
		result.Status = 4003
	} else {
		var posted TSeriesInput
		json.NewDecoder(r.Body).Decode(&posted)
		insertedId := InsertSeries(
			posted.EngName, posted.SeriesName,
		)
		result.Num = insertedId
	}
	result.ResponseWrite(w)
	return nil
}

func UpdateSeriesView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetAdminIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	if userId == "" {
		result.Status = 4003
	} else {
		var posted TSeriesInput
		json.NewDecoder(r.Body).Decode(&posted)
		value := UpdateSeries(posted.EngName, posted.SeriesName, id)
		result.Num = value
	}

	result.ResponseWrite(w)
	return nil
}
