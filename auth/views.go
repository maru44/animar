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

type TRegistForm struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type TTokensForm struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

type TRefreshReturn struct {
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
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

// user info from userId
// url query params(uid)
func SampleGetUserJsonView(w http.ResponseWriter, r *http.Request) error {
	result := helper.TUserJsonResponse{Status: 200}
	query := r.URL.Query()
	uid := query.Get("uid")
	ctx := context.Background()

	user := GetUserFirebase(ctx, uid)
	result.User = *user

	result.ResponseWrite(w)
	return nil
}

// user info from userId
// from cookie
func GetUserModelFCView(w http.ResponseWriter, r *http.Request) error {
	result := helper.TUserJsonResponse{Status: 200}
	userId := helper.GetIdFromCookie(r) //@TODO ここで emailverifiedか確認
	// tokenがキレてたらblankが帰ってくる
	if userId == "" {
		fmt.Print("blank")
		result.Status = 4001
		result.ResponseWrite(w)
		return nil
	}
	ctx := context.Background()

	user := GetUserFirebase(ctx, userId)
	result.User = *user

	result.ResponseWrite(w)
	return nil
}

// login処理
// cookie
func SetJWTCookieView(w http.ResponseWriter, r *http.Request) error {
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
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var tokens TTokensForm
	err = json.Unmarshal(body, &tokens)

	// email か password が間違っていれば blankが帰ってくる
	if tokens.IdToken != "" {
		helper.SetCookiePackage(w, "idToken", tokens.IdToken)
		helper.SetCookiePackage(w, "refreshToken", tokens.RefreshToken)
	} else {
		result.Status = 401
	}
	result.ResponseWrite(w)

	return nil
}

func CreateUserFirstView(w http.ResponseWriter, r *http.Request) error {
	result := helper.TVoidJsonResponse{Status: 200}

	var posted TLoginForm
	json.NewDecoder(r.Body).Decode(&posted)
	posted.ReturnSecureToken = true

	posted_json, _ := json.Marshal(posted)
	url := `https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=` + os.Getenv("FIREBASE_API_KEY")
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
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Print(body)
	result.ResponseWrite(w)

	return nil
}

// refresh idToken
// cookie
func RenewTokenFCView(w http.ResponseWriter, r *http.Request) error {
	result := helper.TVoidJsonResponse{Status: 200}

	// get refresh token from cookie
	refreshToken, _ := r.Cookie("refreshToken")
	if refreshToken.Value == "" {
		result.Status = 402
		return nil
	}

	jsonStr := `{"grant_type": "refresh_token", "refresh_token": "` + refreshToken.Value + `"}`
	url := `https://securetoken.googleapis.com/v1/token?key=` + os.Getenv("FIREBASE_API_KEY")
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tokens TRefreshReturn
	err = json.Unmarshal(body, &tokens)

	if tokens.IdToken != "" {
		helper.DestroyCookie(w, "idToken") // destroy cookie
		helper.SetCookiePackage(w, "idToken", tokens.IdToken)
	} else {
		result.Status = 401
	}

	result.ResponseWrite(w)
	return nil
}

// この流れでclaim取得
// cookie
func TestGetCookie(w http.ResponseWriter, r *http.Request) error {
	result := helper.TVoidJsonResponse{Status: 200}

	claims := helper.GetClaimsFromCookie(r)
	fmt.Print(claims)
	result.ResponseWrite(w)
	return nil
}
