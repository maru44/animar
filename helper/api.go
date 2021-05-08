package helper

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
)

func IsProductionEnv() bool {
	// 本番環境IPリスト
	hosts := []string{
		"aaa",
	}
	host, _ := os.Hostname()

	if runtime.GOOS != "linux" {
		return false
	}
	for _, v := range hosts {
		if v == host {
			return true
		}
	}
	return true
}

type TIntJsonReponse struct {
	Status int `json:"Status"`
	Num    int `json:"ID"`
}

func SetDefaultResponseHeader(w http.ResponseWriter) bool {
	protocol := "http://"
	host := "localhost:3000"
	if IsProductionEnv() {
		protocol = "http://"
		host = "localhost:3000"
	}
	w.Header().Set("Access-Control-Allow-Origin", protocol+host)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Origin, X-Csrftoken, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELTE, PUT")
	return true
}

func (result TIntJsonReponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}
