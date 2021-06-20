package auth

import (
	"animar/v1/configs"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/tools/s3"
	"animar/v1/pkg/tools/tools"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

type TLoginForm struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"` //
	Password    string `json:"password"`
	//PhotoUrl          string `json:"photoUrl"`
	ReturnSecureToken bool `json:"returnSecureToken"`
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

type TCreateReturn struct {
	IdToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
}

type TProfileForm struct {
	DisplayName string `json:"displayName"`
	PhotoUrl    string `json:"photoUrl"`
}

// user info from userId
// url query params(?uid=)
func GetUserModelView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	uid := query.Get("uid")
	ctx := context.Background()

	user := GetUserFirebase(ctx, uid)
	api.JsonResponse(w, map[string]interface{}{"user": user})
	return nil
}

// user info from userId
// from cookie
func GetUserModelFCView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	claims := fire.GetClaimsFromCookie(r)
	// tokenがキレてたらblankが帰ってくる

	switch {
	case userId == "":
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Bad Request")
	case userId != "":
		ctx := context.Background()
		user := GetUserFirebase(ctx, userId)
		api.JsonResponse(w, map[string]interface{}{
			"user":      user,
			"is_verify": claims["email_verified"],
		})
		return nil
	default:
		// if emai is not verified
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Bad Request")
	}
}

// login処理
// cookie
func SetJWTCookieView(w http.ResponseWriter, r *http.Request) error {
	var posted TLoginForm
	json.NewDecoder(r.Body).Decode(&posted)
	posted.ReturnSecureToken = true

	posted_json, _ := json.Marshal(posted)
	url := `https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=` + configs.FirebaseApiKey
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(posted_json),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Bad Request")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var tokens TTokensForm
	err = json.Unmarshal(body, &tokens)

	// email か password が間違っていれば blankが帰ってくる
	if tokens.IdToken != "" {
		api.SetCookiePackage(w, "idToken", tokens.IdToken, 60*60*24)
		api.SetCookiePackage(w, "refreshToken", tokens.RefreshToken, 60*60*24*30)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized")
	}
	api.JsonResponse(w, map[string]interface{}{})
	return nil
}

func CreateUserFirstView(w http.ResponseWriter, r *http.Request) error {
	var posted TLoginForm
	json.NewDecoder(r.Body).Decode(&posted)
	posted.ReturnSecureToken = true
	// if not displayName ===> YYYY@XXXX.XX  >> YYYY
	if posted.DisplayName == "" {
		posted.DisplayName = strings.Split(posted.Email, "@")[0]
	}
	// @TODO 後でちゃんとした画像にする
	// posted.PhotoUrl = fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", configs.Bucket, "ap-northeast-1", "auth/ogp.png")

	posted_json, _ := json.Marshal(posted)
	url := `https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=` + configs.FirebaseApiKey
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(posted_json),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Bad Request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var d TCreateReturn
	err = json.Unmarshal(body, &d)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized")
	}
	api.SetCookiePackage(w, "idToken", d.IdToken, 60*60*24)
	api.SetCookiePackage(w, "refreshToken", d.RefreshToken, 60*60*24*30)

	ctx := context.Background()
	clientAuth := fire.FirebaseClient(ctx)

	// idToken := &d.IdToken
	// defer SetAdminClaim(ctx, clientAuth, *idToken) // set is_admin false
	SendVerifyEmailAtRegister(ctx, clientAuth, posted.Email)

	api.JsonResponse(w, map[string]interface{}{})
	return nil
}

// refresh idToken
// cookie
func RenewTokenFCView(w http.ResponseWriter, r *http.Request) error {
	// get refresh token from cookie
	refreshToken, _ := r.Cookie("refreshToken")
	if refreshToken.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Bad Request")
	}

	jsonStr := `{"grant_type": "refresh_token", "refresh_token": "` + refreshToken.Value + `"}`
	url := `https://securetoken.googleapis.com/v1/token?key=` + configs.FirebaseApiKey
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tools.ErrorLog(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tokens TRefreshReturn
	err = json.Unmarshal(body, &tokens)

	if tokens.IdToken != "" {
		api.DestroyCookie(w, "idToken") // destroy cookie
		api.SetCookiePackage(w, "idToken", tokens.IdToken, 60*60*24)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorize")
	}

	api.JsonResponse(w, map[string]interface{}{})
	return nil
}

// profile 変更
func UserUpdateView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	var posted TProfileForm

	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20MB
	posted.DisplayName = r.FormValue("dname")
	file, fileHeader, err := r.FormFile("image")

	params := (&auth.UserToUpdate{}).
		DisplayName(posted.DisplayName)

	if err == nil {
		defer file.Close()

		returnFileName, err := s3.UploadS3(file, fileHeader.Filename, []string{"auth"})
		if err != nil {
			fmt.Print(err)
		}
		posted.PhotoUrl = returnFileName

		params = (&auth.UserToUpdate{}).
			DisplayName(posted.DisplayName).
			PhotoURL(posted.PhotoUrl)
	}

	ctx := context.Background()
	clientAuth := fire.FirebaseClient(ctx)

	u, err := clientAuth.UpdateUser(ctx, userId, params)
	if err != nil {
		tools.ErrorLog(err)
	}
	api.JsonResponse(w, map[string]interface{}{"user": *u.UserInfo})
	return nil
}

// この流れでclaim取得
// cookie
func TestGetCookie(w http.ResponseWriter, r *http.Request) error {
	claims := fire.GetClaimsFromCookie(r)
	fmt.Print(claims)
	api.JsonResponse(w, map[string]interface{}{})
	return nil
}
