package tools

import (
	"animar/v1/configs"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"firebase.google.com/go/v4/auth"
)

type TUserJsonResponse struct {
	Status     int           `json:"status"`
	User       auth.UserInfo `json:"user"`
	IsVerified bool          `json:"is_verify"`
}

type TVoidJsonResponse struct {
	Status int `json:"status"`
}

type TBaseJsonResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func IsProductionEnv() bool {
	// 本番環境IPリスト
	hosts := strings.Split(configs.ProductionIpList, ", ")
	host, _ := os.Hostname()

	// if runtime.GOOS != "linux" {
	// 	return false
	// }
	for _, v := range hosts {
		if v == host {
			return true
		}
	}
	return false
}

func SetCookiePackage(w http.ResponseWriter, key string, value string) bool {
	var cookie *http.Cookie
	if IsProductionEnv() {
		cookie = &http.Cookie{
			Name:     key,
			Value:    value,
			Path:     "/",
			Domain:   configs.FrontHost,
			MaxAge:   60 * 60 * 24 * 30,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: true,
		}
	} else {
		cookie = &http.Cookie{
			Name:     key,
			Value:    value,
			Path:     "/",
			Domain:   configs.FrontHost,
			MaxAge:   60 * 60 * 24 * 30,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			HttpOnly: true,
		}
	}
	http.SetCookie(w, cookie)
	return true
}

func DestroyCookie(w http.ResponseWriter, key string) bool {
	var cookie *http.Cookie
	if IsProductionEnv() {
		cookie = &http.Cookie{
			Name:     key,
			Value:    "",
			Path:     "",
			Domain:   configs.FrontHost,
			MaxAge:   -1,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: true,
		}
	} else {
		cookie = &http.Cookie{
			Name:     key,
			Value:    "",
			Path:     "",
			Domain:   configs.FrontHost,
			MaxAge:   -1,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			HttpOnly: true,
		}
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
		protocol = "https://"
		host = configs.FrontHost
	}
	w.Header().Set("Access-Control-Allow-Origin", protocol+host)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Origin, X-Csrftoken, Accept, Cookie")
	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	return true
}

func ApiWrapper(fn func()) {
	fn()
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
