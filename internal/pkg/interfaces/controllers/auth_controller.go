package controllers

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/fires"
	"animar/v1/internal/pkg/interfaces/s3"
	"animar/v1/internal/pkg/usecase"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/maru44/perr"
)

type AuthController struct {
	interactor domain.AuthInteractor
	s3         domain.S3Interactor
	CookieController
}

func NewAuthController(firebase fires.Firebase, uploader s3.Uploader) *AuthController {
	return &AuthController{
		interactor: usecase.NewAuthInteractor(
			&fires.AuthRepository{
				Firebase: firebase,
			},
		),
		s3: usecase.NewS3Interactor(
			&s3.S3Repository{
				Uploader: uploader,
			},
		),
		CookieController: *NewCookieController(),
	}
}

func (controller *AuthController) getClaimsFromCookie(r *http.Request) (claims map[string]interface{}, err error) {
	idToken, err := r.Cookie("idToken")
	claims, err = controller.interactor.Claims(idToken.Value)
	return
}

// user model from query ?uid=<userId>
func (controller *AuthController) GetUserModelFromQueryView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uid := query.Get("uid")
	user, err := controller.interactor.UserInfo(uid)

	response(w, r, perr.Wrap(err, perr.Unauthorized), map[string]interface{}{"user": user})
	return
}

// user model from cookie
func (controller *AuthController) GetUserModelFromCookieView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	user, err := controller.interactor.UserInfo(userId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.Unauthorized), nil)
		return
	}

	claims, err := controller.getClaimsFromCookie(r)
	response(w, r, perr.Wrap(err, perr.Unauthorized), map[string]interface{}{"user": user, "is_verify": claims["email_verified"]})
	return
}

// set cookie
func (controller *AuthController) LoginView(w http.ResponseWriter, r *http.Request) {
	var p domain.TLoginForm
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	p.ReturnSecureToken = true
	p_json, err := json.Marshal(p)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	url := `https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=` + configs.FirebaseApiKey
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(p_json),
	)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.InternalServerError), nil)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.InternalServerError), nil)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var tokens domain.TTokensForm
	err = json.Unmarshal(body, &tokens)

	if tokens.IdToken != "" {
		controller.setCookiePackage(w, "idToken", tokens.IdToken, 60*60*24)
		controller.setCookiePackage(w, "refreshToken", tokens.RefreshToken, 60*60*24*30)
	} else {
		response(w, r, perr.New("", perr.Unauthorized), nil)
	}
	response(w, r, perr.Wrap(err, perr.BadRequest), nil)
	return
}

func (controller *AuthController) RegisterView(w http.ResponseWriter, r *http.Request) {
	var p domain.TLoginForm
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	p.ReturnSecureToken = true
	if p.DisplayName == "" {
		p.DisplayName = strings.Split(p.Email, "@")[0]
	}

	p_json, err := json.Marshal(p)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	url := `https://identitytoolkit.googleapis.com/v1/accounts:signUp?key=` + configs.FirebaseApiKey
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(p_json),
	)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.InternalServerError), nil)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.InternalServerError), nil)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var d domain.TCreateReturn
	err = json.Unmarshal(body, &d)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
	} else {
		controller.setCookiePackage(w, "idToken", d.IdToken, 60*60*24)
		controller.setCookiePackage(w, "refreshToken", d.RefreshToken, 60*60*24*30)

		err = controller.interactor.SendVerify(p.Email)
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
	}
	return
}

func (controller *AuthController) RenewTokenView(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie("refreshToken")
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	jsonStr := `{"grant_type": "refresh_token", "refresh_token": "` + refreshToken.Value + `"}`
	url := `https://securetoken.googleapis.com/v1/token?key=` + configs.FirebaseApiKey
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.InternalServerError), nil)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.InternalServerError), nil)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var tokens domain.TRefreshReturn
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	if tokens.IdToken != "" {
		controller.destroyCookie(w, "idToken") // destroy cookie
		controller.setCookiePackage(w, "idToken", tokens.IdToken, 60*60*24)
	} else {
		response(w, r, perr.New("", perr.Unauthorized), nil)
	}
	return
}

func (controller *AuthController) UpdateProfileView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)

	var p domain.TProfileForm
	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20MB
	p.DisplayName = r.FormValue("dname")
	params := domain.TProfileForm{
		DisplayName: p.DisplayName,
	}
	file, fileHeader, err := r.FormFile("image")

	if err == nil {
		defer file.Close()

		returnFileName, err := controller.s3.Image(file, fileHeader.Filename, []string{"auth"})
		if err != nil {
			fmt.Print(err)
		}
		p.PhotoUrl = returnFileName

		params.PhotoUrl = p.PhotoUrl
	}

	user, err := controller.interactor.UpdateProfile(userId, params)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"user": user})
	return
}

func (controller *AuthController) SetJwtTokenView(w http.ResponseWriter, r *http.Request) {
	var p domain.TTokensForm
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	controller.setCookiePackage(w, "idToken", p.IdToken, 60*60*24)
	controller.setCookiePackage(w, "refreshToken", p.RefreshToken, 60*60*24*30)
	response(w, r, err, nil)
	return
}
