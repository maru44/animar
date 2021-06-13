package tools

import (
	"animar/v1/configs"
	"encoding/json"
	"errors"
	"net/http"
)

type TUserToken struct {
	Token string `json:"token,omitempty"`
}

func MethodMiddleware(next http.Handler, methods []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, m := range methods {
			if m == r.Method {
				next.ServeHTTP(w, r)
			}
		}
		http.Error(w, http.StatusText(405), 405)
		return
	})
}

func UpsertOnlyMiddleware(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	http.Error(w, "METHOD NOT ALLOWED", 405)
	return errors.New("METHOD NOT ALLOWED")
}

func PostOnlyMiddleware(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" || r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	http.Error(w, "METHOD NOT ALLOWED", 405)
	return errors.New("METHOD NOT ALLOWED")
}

func PutOnlyMiddleware(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "PUT" || r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	http.Error(w, "METHOD NOT ALLOWED", 405)
	return errors.New("METHOD NOT ALLOWED")
}

func DeleteOnlyMiddleware(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "DELETE" || r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return nil
	}
	http.Error(w, "METHOD NOT ALLOWED", 405)
	return errors.New("METHOD NOT ALLOWED")
}

// for CSR
func AdminRequiredMiddleware(w http.ResponseWriter, r *http.Request) error {
	userId := GetAdminIdFromCookie(r)
	isAdmin := IsAdmin(userId)
	if !isAdmin {
		http.Error(w, "FORBIDDEN", 403)
		return errors.New("FORBIDDEN")
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// for SSR get & CSR get
func AdminRequiredMiddlewareGet(w http.ResponseWriter, r *http.Request) error {
	var userId string
	switch r.Method {
	case "GET":
		userId = GetAdminIdFromCookie(r)
	case "POST":
		var p TUserToken
		json.NewDecoder(r.Body).Decode(&p)
		userId = GetAdminIdFromIdToken(p.Token)
	default:
		userId = ""
	}
	if userId == "" {
		http.Error(w, "FORBIDDEN", 403)
		return errors.New("FORBIDDEN")
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

func corsMiddleware(w http.ResponseWriter, r *http.Request) error {
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
	return nil
}
