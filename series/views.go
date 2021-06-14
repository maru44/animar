package series

import (
	"animar/v1/tools/api"
	"animar/v1/tools/fire"
	"encoding/json"
	"net/http"
	"strconv"
)

func SeriesView(w http.ResponseWriter, r *http.Request) error {
	result := api.TBaseJsonResponse{Status: 200}
	var userId string

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
	result := api.TBaseJsonResponse{Status: 200}

	var p TSeriesInput
	json.NewDecoder(r.Body).Decode(&p)
	insertedId := InsertSeries(
		p.EngName, p.SeriesName,
	)
	result.Data = insertedId

	result.ResponseWrite(w)
	return nil
}

func UpdateSeriesView(w http.ResponseWriter, r *http.Request) error {
	result := api.TBaseJsonResponse{Status: 200}

	result, is_valid := result.LimitMethod([]string{"PUT"}, r)
	if !is_valid {
		result.ResponseWrite(w)
		return nil
	}

	userId := fire.GetAdminIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	if userId == "" {
		result.Status = 4003
	} else {
		var posted TSeriesInput
		json.NewDecoder(r.Body).Decode(&posted)
		value := UpdateSeries(posted.EngName, posted.SeriesName, id)
		result.Data = value
	}

	result.ResponseWrite(w)
	return nil
}
