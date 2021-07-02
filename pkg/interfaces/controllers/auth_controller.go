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
	}
}

func (controller *AuthController) GetUserFromQueryView(w http.ResponseWriter, r *http.Request) (ret error) {
	userId := r.URL.Query().Get("uid")

	user, err := controller.interactor.UserInfo(userId)
	ret = response(w, err, map[string]interface{}{"user": user})
	return
}

func (controller *AuthController) GetUserFromCookieView(w http.ResponseWriter, r *http.Request) (ret error) {
	claims, err := controller.getClaimsFromCookie(r)
	userId := claims["user_id"].(string)
	user, err := controller.interactor.UserInfo(userId)
	ret = response(w, err, map[string]interface{}{"user": user, "is_verify": claims["email_verified"]})
	return
}

// set cookie
func (controller *AuthController) LoginView(w http.ResponseWriter, r *http.Request) (ret error) {
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
		ret = response(w, err, nil)
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
		ret = response(w, domain.ErrUnauthorized, nil)
	}
	ret = response(w, err, nil)
	return ret
}

func (controller *AuthController) RegisterView(w http.ResponseWriter, r *http.Request) (ret error) {
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
		ret = response(w, err, nil)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var d domain.TCreateReturn
	err = json.Unmarshal(body, &d)
	if err != nil {
		ret = response(w, err, nil)
	} else {
		api.SetCookiePackage(w, "idToken", d.IdToken, 60*60*24)
		api.SetCookiePackage(w, "refreshToken", d.RefreshToken, 60*60*24*30)

		err = controller.interactor.SendVerify(p.Email)
		ret = response(w, err, nil)
	}
	return ret
}

func (controller *AuthController) RenewTokenView(w http.ResponseWriter, r *http.Request) error {
	refreshToken, _ := r.Cookie("refreshToken")
	if refreshToken.Value == "" {
		return response(w, domain.ErrBadRequest, nil)
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
		return response(w, domain.ErrUnauthorized, nil)
	}
	return nil
}

func (controller *AuthController) UpdateProfileView(w http.ResponseWriter, r *http.Request) (ret error) {
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
	ret = response(w, err, map[string]interface{}{"user": user})
	return ret
}

/************************
        middleware
*************************/

func (controller *AuthController) AdminRequiredMiddleware(w http.ResponseWriter, r *http.Request) (ret error) {
	idToken, err := r.Cookie("idToken")
	if err != nil {
		ret = response(w, err, nil)
		return
	}
	_, err = controller.interactor.AdminId(idToken.Value)
	if err != nil {
		ret = response(w, err, nil)
	} else {
		err = nil
	}
	return
}

func (controller *AuthController) SSRAdminRequiredMiddleware(w http.ResponseWriter, r *http.Request) error {
	var idToken string
	var err error
	switch r.Method {
	case "GET":
		token, err := r.Cookie("idToken")
		if err != nil {
			return response(w, err, nil)
		}
		idToken = token.Value
	case "POST":
		var p domain.TUserToken
		json.NewDecoder(r.Body).Decode(&p)
		idToken = p.Token
	default:
		err = domain.ErrForbidden
		return err
	}
	_, err = controller.interactor.AdminId(idToken)
	return err
}
