package tools

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"

	"firebase.google.com/go/v4/auth"
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

type TStringJsonResponse struct {
	Status int    `json:"Status"`
	String string `json:"String"`
}

type TUserJsonResponse struct {
	Status     int           `json:"Status"`
	User       auth.UserInfo `json:"User"`
	IsVerified bool          `json:"IsVerify"`
}

type TVoidJsonResponse struct {
	Status int `json:"Status"`
}

type TBaseJsonResponse struct {
	Status int         `json:"Status"`
	Data   interface{} `json:"Data"`
}

// @TODO env使う
func SetCookiePackage(w http.ResponseWriter, key string, value string) bool {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   60 * 60 * 24 * 30,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return true
}

// @TODO env使う
func DestroyCookie(w http.ResponseWriter, key string) bool {
	cookie := &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "",
		Domain:   "localhost",
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	return true
}

type Api interface {
	ResponseWrite(w http.ResponseWriter) bool
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Origin, X-Csrftoken, Accept, Cookie")
	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
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

func (result TUserJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
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

func (result TVoidJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
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

func (result TBaseJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
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

func (result TStringJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
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
