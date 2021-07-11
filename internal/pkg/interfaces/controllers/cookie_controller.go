package controllers

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/tools/tools"
	"net/http"
)

type CookieController struct{}

func NewCookieController() *CookieController {
	return &CookieController{}
}

func (c *CookieController) setCookiePackage(w http.ResponseWriter, key string, value string, age int) bool {
	var cookie *http.Cookie
	if tools.IsProductionEnv() {
		cookie = &http.Cookie{
			Name:     key,
			Value:    value,
			Path:     "/",
			Domain:   configs.FrontHost,
			MaxAge:   age,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: true,
		}
	} else {
		cookie = &http.Cookie{
			Name:     key,
			Value:    value,
			Path:     "/",
			Domain:   configs.FrontHost,
			MaxAge:   age,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			HttpOnly: true,
		}
	}
	http.SetCookie(w, cookie)
	return true
}

func (c *CookieController) destroyCookie(w http.ResponseWriter, key string) bool {
	var cookie *http.Cookie
	if tools.IsProductionEnv() {
		cookie = &http.Cookie{
			Name:     key,
			Value:    "",
			Path:     "",
			Domain:   configs.FrontHost,
			MaxAge:   -1,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: true,
		}
	} else {
		cookie = &http.Cookie{
			Name:     key,
			Value:    "",
			Path:     "",
			Domain:   configs.FrontHost,
			MaxAge:   -1,
			SameSite: http.SameSiteLaxMode,
			Secure:   false,
			HttpOnly: true,
		}
	}
	http.SetCookie(w, cookie)
	return true
}
