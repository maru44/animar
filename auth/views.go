package auth

import (
	"animar/v1/helper"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type TLoginForm struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type TTokensForm struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

func UserListView(w http.ResponseWriter, r *http.Request) {
	jsonStr := `{"grant_type": "client_credentials", "client_id": "` + os.Getenv("AUTH0_CLIENT_ID") + `", "client_secret": "` + os.Getenv("AUTH0_SECRET") + `", "audience": "https://` + os.Getenv("AUTH0_DOMAIN") + `/api/v2/"}`
	url := `https://` + os.Getenv("AUTH0_DOMAIN") + `/oauth/token`
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	fmt.Fprintln(w, resp)
}

func SampleGetUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uid := query.Get("uid")
	ctx := context.Background()

	user := GetUserFirebase(ctx, uid)
	fmt.Fprintln(w, user)
}

func SampleGetUserJson(w http.ResponseWriter, r *http.Request) error {
	result := helper.TUserJsonResponse{Status: 200}
	query := r.URL.Query()
	uid := query.Get("uid")
	ctx := context.Background()

	user := GetUserFirebase(ctx, uid)
	result.User = *user

	helper.SetCookiePackage(w, "aaa", "bbbbb")
	result.ResponseWrite(w)
	return nil
}

func SetJWTCookie(w http.ResponseWriter, r *http.Request) error {
	result := helper.TVoidJsonResponse{Status: 200}

	var posted TLoginForm
	json.NewDecoder(r.Body).Decode(&posted)
	posted.ReturnSecureToken = true

	posted_json, _ := json.Marshal(posted)
	url := `https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=` + os.Getenv("FIREBASE_API_KEY")
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(posted_json),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Status = 400
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var tokens TTokensForm
	err = json.Unmarshal(body, &tokens)

	helper.SetCookiePackage(w, "idToken", tokens.IdToken)
	helper.SetCookiePackage(w, "refreshToken", tokens.RefreshToken)
	result.ResponseWrite(w)

	return nil
}
