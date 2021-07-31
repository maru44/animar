package controllers

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/infrastructure"
	"animar/v1/internal/pkg/interfaces/fires"
	"animar/v1/internal/pkg/tools/tools"
	"animar/v1/internal/pkg/usecase"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type BaseController struct {
	interactor domain.BaseInteractor
}

type ContextKey string

const (
	USER_ID ContextKey = "userId"
)

func NewBaseController() *BaseController {
	return &BaseController{
		interactor: usecase.NewBaseInteractor(
			&fires.AuthRepository{
				Firebase: infrastructure.NewFireBaseClient(),
			},
		),
	}
}

/************************
         view
************************/

func (controller *BaseController) GatewayView(w http.ResponseWriter, r *http.Request) {
	response(w, domain.ErrNotFound, nil)
	return
}

/************************
         user
************************/

func (controller *BaseController) getUserIdFromCookie(r *http.Request) (userId string, err error) {
	idToken, err := r.Cookie("idToken")
	if err != nil {
		return
	} else if idToken.Value == "" {
		return
	}
	userId, err = controller.interactor.UserId(idToken.Value)
	return
}

func (controller *BaseController) getUserIdFromToken(idToken string) (userId string, err error) {
	claims, err := controller.interactor.Claims(idToken)
	if err != nil {
		return
	}
	userId = claims["user_id"].(string)
	return
}

func (controller *BaseController) getGoogleUser(accessToken string) domain.TGoogleOauth {
	url := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken
	req, err := http.NewRequest(
		"GET", url, nil,
	)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		domain.ErrorWarn(err)
	}
	defer res.Body.Close()

	defer res.Body.Close()
	byteArray, _ := ioutil.ReadAll(res.Body)
	var user domain.TGoogleOauth
	err = json.Unmarshal(byteArray, &user)
	return user
}

/************************
    method middleware
************************/

func (controller *BaseController) allowOptionsMiddleware(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		response(w, nil, nil)
	}
}

func (controller *BaseController) UpsertOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" {
			next.ServeHTTP(w, r)
		} else {
			response(w, domain.ErrMethodNotAllowed, nil)
			return
		}
	})
}

func (controller *BaseController) PostOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			next.ServeHTTP(w, r)
		} else {
			response(w, domain.ErrMethodNotAllowed, nil)
			return
		}
	})
}

func (controller *BaseController) PutOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PUT" {
			next.ServeHTTP(w, r)
		} else {
			response(w, domain.ErrMethodNotAllowed, nil)
			return
		}
	})
}

func (controller *BaseController) DeleteOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			next.ServeHTTP(w, r)
		} else {
			response(w, domain.ErrMethodNotAllowed, nil)
			return
		}
	})
}

/************************
        middleware
************************/

func (controller *BaseController) corsMiddleware(w http.ResponseWriter, r *http.Request) error {
	protocol := "http://"
	host := "localhost:3000"
	if tools.IsProductionEnv() {
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

func (controller *BaseController) BaseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller.corsMiddleware(w, r)
		controller.allowOptionsMiddleware(w, r)
		next.ServeHTTP(w, r)
	})
}

/************************
    user middleware
************************/

// set context userId
func (controller *BaseController) GiveUserIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken, err := r.Cookie("idToken")
		var userId string
		if err != nil {
			userId = ""
		} else {
			userId, err = controller.interactor.UserId(idToken.Value)
			if err != nil {
				userId = ""
			}
		}
		ctx := context.WithValue(r.Context(), USER_ID, userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// set context userId (SSR and CSR)
func (controller *BaseController) GiveUserIdMiddlewareAbleSSR(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId string
		switch r.Method {
		case "GET":
			idToken, err := r.Cookie("idToken")
			if err != nil {
				userId = ""
			} else {
				userId, _ = controller.interactor.UserId(idToken.Value)
			}
		case "POST":
			var posted domain.TUserToken
			json.NewDecoder(r.Body).Decode(&posted)
			userId, _ = controller.getUserIdFromToken(posted.Token)
		default:
			userId = ""
		}
		ctx := context.WithValue(r.Context(), USER_ID, userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// must login and set context userId
func (controller *BaseController) LoginRequireMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken, err := r.Cookie("idToken")
		if err != nil {

			response(w, domain.ErrUnauthorized, nil)
			return
		}
		userId, err := controller.interactor.UserId(idToken.Value)
		if err != nil {
			response(w, domain.ErrUnauthorized, nil)
			return
		}
		ctx := context.WithValue(r.Context(), USER_ID, userId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

/************************
        admin
************************/

// for CSR
func (controller *BaseController) AdminRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken, err := r.Cookie("idToken")
		if err != nil {
			response(w, domain.ErrForbidden, nil)
			return
		} else if idToken.Value == "" {
			response(w, domain.ErrForbidden, nil)
			return
		}
		userId, err := controller.interactor.AdminId(idToken.Value)
		if err != nil {
			response(w, domain.ErrForbidden, nil)
			return
		} else if userId == "" {
			response(w, domain.ErrForbidden, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// for SSR get & CSR get
func (controller *BaseController) AdminRequiredMiddlewareGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var idToken string
		switch r.Method {
		case "GET":
			idTokenObj, err := r.Cookie("idToken")
			if err != nil {
				response(w, domain.ErrForbidden, nil)
				return
			} else {
				idToken = idTokenObj.Value
			}
		case "POST":
			var p domain.TUserToken
			json.NewDecoder(r.Body).Decode(&p)
			idToken = p.Token
		default:
			response(w, domain.ErrForbidden, nil)
			return
		}

		userId, err := controller.interactor.AdminId(idToken)
		if err != nil {
			response(w, domain.ErrForbidden, nil)
			return
		} else if userId == "" {
			response(w, domain.ErrForbidden, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
