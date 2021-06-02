package series

import (
	"animar/v1/tools"
	"encoding/json"
	"net/http"
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

func ListSeriesView(w http.ResponseWriter, r *http.Request) error {
	result := TSeriesResponse{Status: 200}

	var sers []TSeries
	sers = ListSeriesDomain()

	result.Data = sers
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
