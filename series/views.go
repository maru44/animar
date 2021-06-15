package series

import (
	"animar/v1/tools/api"
	"animar/v1/tools/fire"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func SeriesView(w http.ResponseWriter, r *http.Request) error {
	var userId string

	query := r.URL.Query()
	id := query.Get("id")

	if userId == "" {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("Forbidden")
	} else {
		if id != "" {
			i, _ := strconv.Atoi(id)
			ser := DetailSeries(i)
			if ser.ID == 0 {
				w.WriteHeader(http.StatusNotFound)
				return errors.New("Not Found")
			}
			api.JsonResponse(w, map[string]interface{}{"data": ser})
		} else {
			sers := ListSeriesDomain()
			api.JsonResponse(w, map[string]interface{}{"data": sers})
		}
	}
	return nil
}

func InsertSeriesView(w http.ResponseWriter, r *http.Request) error {
	var p TSeriesInput
	json.NewDecoder(r.Body).Decode(&p)
	insertedId := InsertSeries(
		p.EngName, p.SeriesName,
	)
	api.JsonResponse(w, map[string]interface{}{"data": insertedId})
	return nil
}

func UpdateSeriesView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetAdminIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	if userId == "" {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("Forbidden")
	} else {
		var posted TSeriesInput
		json.NewDecoder(r.Body).Decode(&posted)
		value := UpdateSeries(posted.EngName, posted.SeriesName, id)
		api.JsonResponse(w, map[string]interface{}{"data": value})
	}
	return nil
}
