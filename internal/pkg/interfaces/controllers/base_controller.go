package controllers

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/infrastructure"
	"animar/v1/internal/pkg/interfaces/fires"
	"animar/v1/internal/pkg/tools/tools"
	"animar/v1/internal/pkg/usecase"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/maru44/perr"
)

type BaseController struct {
	interactor domain.BaseInteractor
	cache      domain.Cache
	CookieController
}

const (
	CSRF_COOKIE_KEY = "csrf-token"
)

func NewBaseController(c domain.Cache) *BaseController {
	return &BaseController{
		interactor: usecase.NewBaseInteractor(
			&fires.AuthRepository{
				Firebase: infrastructure.NewFireBaseClient(),
			},
		),
		cache: c,
	}
}

/************************
         view
************************/

func (controller *BaseController) GatewayView(w http.ResponseWriter, r *http.Request) {
	response(w, r, perr.New("", perr.NotFound), nil)
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
	req, _ := http.NewRequest(
		http.MethodGet, url, nil,
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
	json.Unmarshal(byteArray, &user)
	return user
}

/************************
    method middleware
************************/

func (controller *BaseController) UpsertOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			next.ServeHTTP(w, r)
		} else if r.Method == http.MethodOptions {
			response(w, r, nil, nil)
		} else {
			response(w, r, perr.New("", perr.MethodNotAllowed), nil)
			return
		}
	})
}

func (controller *BaseController) PostOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
		} else if r.Method == http.MethodOptions {
			response(w, r, nil, nil)
		} else {
			response(w, r, perr.New("", perr.MethodNotAllowed), nil)
			return
		}
	})
}

func (controller *BaseController) PutOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			next.ServeHTTP(w, r)
		} else if r.Method == http.MethodOptions {
			response(w, r, nil, nil)
		} else {
			response(w, r, perr.New("", perr.MethodNotAllowed), nil)
			return
		}
	})
}

func (controller *BaseController) DeleteOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			next.ServeHTTP(w, r)
		} else if r.Method == http.MethodOptions {
			response(w, r, nil, nil)
		} else {
			response(w, r, perr.New("", perr.MethodNotAllowed), nil)
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
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Origin, X-Csrftoken, Accept, Cookie")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	w.Header().Set("Access-Control-Max-Age", "3600")
	return nil
}

func (controller *BaseController) BaseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = setAccessData(r)

		controller.corsMiddleware(w, r)
		next.ServeHTTP(w, r)
	})
}

// csrf
func (controller *BaseController) VerifyCsrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfToken, err := r.Cookie(CSRF_COOKIE_KEY)
		if err != nil {
			err = perr.Wrap(err, perr.BadRequest, "Invalid Csrf Token")
			response(w, r, err, nil)
			return
		}
		cache := controller.cache
		if ok := cache.Get(domain.CacheTypeCsrf, csrfToken.Value); !ok {
			err = perr.New("Invalid Csrf Token", perr.BadRequest)
			response(w, r, err, nil)
			return
		}
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
		r = setUserToContext(r, userId)
		next.ServeHTTP(w, r)
	})
}

// set context userId (SSR and CSR)
func (controller *BaseController) GiveUserIdMiddlewareAbleSSR(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId string
		switch r.Method {
		case http.MethodGet:
			idToken, err := r.Cookie("idToken")
			if err != nil {
				userId = ""
			} else {
				userId, _ = controller.interactor.UserId(idToken.Value)
			}
		case http.MethodPost:
			var posted domain.TUserToken
			json.NewDecoder(r.Body).Decode(&posted)
			userId, _ = controller.getUserIdFromToken(posted.Token)
		default:
			userId = ""
		}
		r = setUserToContext(r, userId)
		next.ServeHTTP(w, r)
	})
}

// must login and set context userId
func (controller *BaseController) LoginRequireMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken, err := r.Cookie("idToken")
		if err != nil {
			response(w, r, perr.Wrap(err, perr.Expired), nil)
			return
		}
		userId, err := controller.interactor.UserId(idToken.Value)
		if err != nil {
			response(w, r, perr.Wrap(err, perr.Unauthorized), nil)
			return
		}
		r = setUserToContext(r, userId)
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
			response(w, r, perr.Wrap(err, perr.Forbidden), nil)
			return
		} else if idToken.Value == "" {
			response(w, r, perr.New("", perr.Unauthorized), nil)
			return
		}
		userId, err := controller.interactor.AdminId(idToken.Value)
		if err != nil {
			response(w, r, perr.Wrap(err, perr.Unauthorized), nil)
			return
		} else if userId == "" {
			response(w, r, perr.New("", perr.Unauthorized), nil)
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
		case http.MethodGet:
			idTokenObj, err := r.Cookie("idToken")
			if err != nil {
				response(w, r, perr.Wrap(err, perr.Forbidden), nil)
				return
			} else {
				idToken = idTokenObj.Value
			}
		case http.MethodPost:
			var p domain.TUserToken
			json.NewDecoder(r.Body).Decode(&p)
			idToken = p.Token
		default:
			response(w, r, perr.New("", perr.Forbidden), nil)
			return
		}

		userId, err := controller.interactor.AdminId(idToken)
		if err != nil {
			response(w, r, err, nil)
			return
		} else if userId == "" {
			response(w, r, perr.New("", perr.Forbidden), nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

/************************
      set cookie
************************/

func (controller *BaseController) SetCsrfCookieView(w http.ResponseWriter, r *http.Request) {
	csrfToken := tools.GenRandSlug(32)
	controller.cache.AddCacheItem(domain.CacheTypeCsrf, csrfToken, domain.CSRF_INTERVAL_MINUTE*time.Minute)
	controller.destroyCookie(w, CSRF_COOKIE_KEY)
	controller.setCookiePackage(w, CSRF_COOKIE_KEY, csrfToken, 60*domain.CSRF_INTERVAL_MINUTE)
	return
}
