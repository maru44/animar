package controllers

import (
	"animar/v1/configs"
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/fires"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/s3"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthController struct {
	interactor domain.AuthInteractor
	BaseController
}

func NewAuthController(firebase fires.Firebase) *AuthController {
	return &AuthController{
		interactor: usecase.NewAuthInteractor(
			&fires.AuthRepository{
				Firebase: firebase,
			},
		),
		BaseController: *NewBaseController(),
		// BaseController: BaseController{
		// 	interactor: usecase.NewBaseInteractor(
		// 		&fires.AuthRepository{
		// 			Firebase: firebase,
		// 		},
		// 	),
		// },
	}
}

// user modal from query ?uid=<userId>
func (controller *AuthController) GetUserModelFromQueryView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uid := query.Get("uid")
	user, err := controller.interactor.UserInfo(uid)

	response(w, err, map[string]interface{}{"user": user})
	return
}

// user model from cookie
func (controller *AuthController) GetUserModelFromCookieView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	user, err := controller.interactor.UserInfo(userId)
	claims, err := controller.getClaimsFromCookie(r)
	response(w, err, map[string]interface{}{"user": user, "is_verify": claims["email_verified"]})
	return
}

// set cookie
func (controller *AuthController) LoginView(w http.ResponseWriter, r *http.Request) {
	var p domain.TLoginForm
	json.NewDecoder(r.Body).Decode(&p)
	p.ReturnSecureToken = true

	p_json, _ := json.Marshal(p)
	url := `https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=` + configs.FirebaseApiKey
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(p_json),
	)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		response(w, err, nil)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var tokens domain.TTokensForm
	err = json.Unmarshal(body, &tokens)

	if tokens.IdToken != "" {
		api.SetCookiePackage(w, "idToken", tokens.IdToken, 60*60*24)
		api.SetCookiePackage(w, "refreshToken", tokens.RefreshToken, 60*60*24*30)
	} else {
		response(w, domain.ErrUnauthorized, nil)
	}
	response(w, err, nil)
	return
}

func (controller *AuthController) RegisterView(w http.ResponseWriter, r *http.Request) {
	var p domain.TLoginForm
	json.NewDecoder(r.Body).Decode(&p)
	p.ReturnSecureToken = true

	if p.DisplayName == "" {
		p.DisplayName = strings.Split(p.Email, "@")[0]
	}
	p_json, _ := json.Marshal(p)
	url := `https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=` + configs.FirebaseApiKey
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(p_json),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		response(w, err, nil)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var d domain.TCreateReturn
	err = json.Unmarshal(body, &d)
	if err != nil {
		response(w, err, nil)
	} else {
		api.SetCookiePackage(w, "idToken", d.IdToken, 60*60*24)
		api.SetCookiePackage(w, "refreshToken", d.RefreshToken, 60*60*24*30)

		err = controller.interactor.SendVerify(p.Email)
		response(w, err, nil)
	}
	return
}

func (controller *AuthController) RenewTokenView(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		response(w, domain.ErrUnauthorized, nil)
		return
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
	var tokens domain.TRefreshReturn
	err = json.Unmarshal(body, &tokens)

	if tokens.IdToken != "" {
		api.DestroyCookie(w, "idToken") // destroy cookie
		api.SetCookiePackage(w, "idToken", tokens.IdToken, 60*60*24)
	} else {
		response(w, domain.ErrUnauthorized, nil)
	}
	return
}

func (controller *AuthController) UpdateProfileView(w http.ResponseWriter, r *http.Request) {
	claims, err := controller.getClaimsFromCookie(r)
	userId := claims["user_id"].(string)

	var p domain.TProfileForm
	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20MB
	p.DisplayName = r.FormValue("dname")
	file, fileHeader, err := r.FormFile("image")

	params := domain.TProfileForm{
		DisplayName: p.DisplayName,
	}

	if err == nil {
		defer file.Close()

		returnFileName, err := s3.UploadS3(file, fileHeader.Filename, []string{"auth"})
		if err != nil {
			fmt.Print(err)
		}
		p.PhotoUrl = returnFileName

		params.PhotoUrl = p.PhotoUrl
	}

	user, err := controller.interactor.UpdateProfile(userId, params)
	response(w, err, map[string]interface{}{"user": user})
	return
}

func (controller *AuthController) GoogleOAuthView(w http.ResponseWriter, r *http.Request) {
	controller.interactor.OauthGoogle()
	response(w, nil, nil)
	return
}

func (controller *AuthController) GoogleRedirectView(w http.ResponseWriter, r *http.Request) {
	controller.interactor.GoogleRedirect(r.FormValue("code"))
	response(w, nil, nil)
	return
}

// func (controller *AuthController) CreateGoogleLinkView(w http.ResponseWriter, r *http.Request) {
// 	state := "aaa"
// 	u, err := url.Parse(authorization)
// }
