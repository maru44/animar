package api

import (
	"animar/v1/configs"
	"animar/v1/pkg/tools/tools"
	"encoding/json"
	"net/http"
)

type TBaseJsonResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

type TUserToken struct {
	Token string `json:"token,omitempty"`
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

func JsonResponse(w http.ResponseWriter, dictionary map[string]interface{}) bool {
	data, err := json.Marshal(dictionary)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return true
}
