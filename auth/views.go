package auth

import (
	"animar/v1/tools"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	IdToken string `json:"idToken"`
	Email   string `json:"email"`
}

type TProfileForm struct {
	DisplayName string `json:"displayName"`
	PhotoUrl    string `json:"photoUrl"`
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
		tools.ErrorLog(err)
	}
	defer resp.Body.Close()

	fmt.Fprintln(w, resp)
}

// user info from userId
// url query params(?uid=)
func GetUserModelView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TUserJsonResponse{Status: 200}
	query := r.URL.Query()
	uid := query.Get("uid")
	ctx := context.Background()

	user := GetUserFirebase(ctx, uid)
	if user != nil {
		result.User = *user
	}

	result.ResponseWrite(w)
	return nil
}

// user info from userId
// from cookie
func GetUserModelFCView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TUserJsonResponse{Status: 200}
	userId := tools.GetIdFromCookie(r)
	claims := tools.GetClaimsFromCookie(r)
	// tokenがキレてたらblankが帰ってくる

	switch {
	case userId == "":
		result.Status = 4001
		/*
			// メール未承認入れない場合
			case claims["email_verified"]:
				ctx := context.Background()
				user := GetUserFirebase(ctx, userId)
				result.User = *user
		*/
	case userId != "":
		ctx := context.Background()
		user := GetUserFirebase(ctx, userId)
		result.User = *user
		if claims["email_verified"] == true {
			result.IsVerified = true
		} else {
			result.IsVerified = false
		}
	default:
		// if emai is not verified
		result.Status = 4002
	}

	result.ResponseWrite(w)
	return nil
}

// login処理
// cookie
func SetJWTCookieView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TVoidJsonResponse{Status: 200}

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
		tools.SetCookiePackage(w, "idToken", tokens.IdToken)
		tools.SetCookiePackage(w, "refreshToken", tokens.RefreshToken)
	} else {
		result.Status = 401
	}
	result.ResponseWrite(w)

	return nil
}

func CreateUserFirstView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TVoidJsonResponse{Status: 200}

	var posted TLoginForm
	json.NewDecoder(r.Body).Decode(&posted)
	posted.ReturnSecureToken = true
	// if not displayName ===> YYYY@XXXX.XX  >> YYYY
	if posted.DisplayName == "" {
		posted.DisplayName = strings.Split(posted.Email, "@")[0]
	}
	// @TODO 後でちゃんとした画像にする
	// posted.PhotoUrl = fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", os.Getenv("BUCKET"), "ap-northeast-1", "auth/ogp.png")

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
	body, err := ioutil.ReadAll(resp.Body)

	var d TCreateReturn
	err = json.Unmarshal(body, &d)
	if err != nil {
		result.Status = 400
		return nil
	}
	tools.SetCookiePackage(w, "idToken", d.IdToken)

	ctx := context.Background()
	clientAuth := tools.FirebaseClient(ctx)

	// idToken := &d.IdToken
	// defer SetAdminClaim(ctx, clientAuth, *idToken) // set is_admin false
	SendVerifyEmailAtRegister(ctx, clientAuth, posted.Email)

	result.ResponseWrite(w)
	return nil
}

// refresh idToken
// cookie
func RenewTokenFCView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TVoidJsonResponse{Status: 200}

	// get refresh token from cookie
	refreshToken, _ := r.Cookie("refreshToken")
	if refreshToken.Value == "" {
		result.Status = 4002
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
		tools.ErrorLog(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tokens TRefreshReturn
	err = json.Unmarshal(body, &tokens)

	if tokens.IdToken != "" {
		tools.DestroyCookie(w, "idToken") // destroy cookie
		tools.SetCookiePackage(w, "idToken", tokens.IdToken)
	} else {
		result.Status = 401
	}

	result.ResponseWrite(w)
	return nil
}

// profile 変更
func UserUpdateView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TUserJsonResponse{Status: 200}

	userId := tools.GetIdFromCookie(r)
	var posted TProfileForm

	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20MB
	posted.DisplayName = r.FormValue("dname")
	file, fileHeader, err := r.FormFile("image")

	params := (&auth.UserToUpdate{}).
		DisplayName(posted.DisplayName)

	if err == nil {
		defer file.Close()

		returnFileName, err := tools.UploadS3(file, fileHeader.Filename, []string{"auth"})
		if err != nil {
			fmt.Print(err)
		}
		posted.PhotoUrl = returnFileName

		params = (&auth.UserToUpdate{}).
			DisplayName(posted.DisplayName).
			PhotoURL(posted.PhotoUrl)
	}

	ctx := context.Background()
	clientAuth := tools.FirebaseClient(ctx)

	u, err := clientAuth.UpdateUser(ctx, userId, params)
	if err != nil {
		tools.ErrorLog(err)
	}
	result.User = *u.UserInfo

	result.ResponseWrite(w)
	return nil
}

// この流れでclaim取得
// cookie
func TestGetCookie(w http.ResponseWriter, r *http.Request) error {
	result := tools.TVoidJsonResponse{Status: 200}

	claims := tools.GetClaimsFromCookie(r)
	fmt.Print(claims)
	result.ResponseWrite(w)
	return nil
}
