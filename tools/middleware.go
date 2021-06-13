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

func allowOptionsMiddleware(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "OPTIONS" {
		//w.WriteHeader(http.StatusOK)
		return nil
	}
	return nil
}

func UpsertOnlyMiddleware(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" || r.Method == "PUT" {
		return nil
	}
	http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)
	return errors.New("METHOD NOT ALLOWED")
}

func PostOnlyMiddleware(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return nil
	}
	http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)
	return errors.New("METHOD NOT ALLOWED")
}

func PutOnlyMiddleware(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "PUT" {
		return nil
	}
	http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)
	return errors.New("METHOD NOT ALLOWED")
}

func DeleteOnlyMiddleware(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "DELETE" {
		return nil
	}
	http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)
	return errors.New("METHOD NOT ALLOWED")
}

// for CSR
func AdminRequiredMiddleware(w http.ResponseWriter, r *http.Request) error {
	userId := GetAdminIdFromCookie(r)
	isAdmin := IsAdmin(userId)
	if !isAdmin {
		http.Error(w, "FORBIDDEN", http.StatusForbidden)
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
		http.Error(w, "FORBIDDEN", http.StatusForbidden)
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
	//w.Header().Set("Content-Type", "application/json, multipart/formdata, text/plain")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Origin, X-Csrftoken, Accept, Cookie")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	w.Header().Set("Access-Control-Max-Age", "3600")
	return nil
}
