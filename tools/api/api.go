package api

import (
	"animar/v1/configs"
	"animar/v1/tools/tools"
	"encoding/json"
	"net/http"

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

func SetCookiePackage(w http.ResponseWriter, key string, value string, age int) bool {
	var cookie *http.Cookie
	if tools.IsProductionEnv() {
		cookie = &http.Cookie{
			Name:     key,
			Value:    value,
			Path:     "/",
			Domain:   configs.FrontHost,
			MaxAge:   age,
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
			MaxAge:   age,
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
	if tools.IsProductionEnv() {
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

func (result TUserJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	w.Write(res)
	w.WriteHeader(http.StatusOK)
	return true
}

func (result TVoidJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	w.Write(res)
	w.WriteHeader(http.StatusOK)
	return true
}

func (result TBaseJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	w.Write(res)
	w.WriteHeader(http.StatusOK)
	return true
}

func (result TBaseJsonResponse) LimitMethod(validMethods []string, r *http.Request) (TBaseJsonResponse, bool) {
	method := r.Method
	for _, m := range validMethods {
		if method == m {
			return result, true
		}
	}
	result.Status = 4005
	return result, false
}
