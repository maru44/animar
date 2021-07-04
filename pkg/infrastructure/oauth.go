package infrastructure

import (
	"animar/v1/pkg/interfaces/fires"

	"golang.org/x/oauth2"
)

type Oauth struct {
	Conn *oauth2.Config
}

type OauthToken struct {
	Token *oauth2.Token
}

func NewGoogleConfig() fires.Oauth {
	conf := &oauth2.Config{
		ClientID:     "googleClientID",
		ClientSecret: "googleClientSecret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "aaaa",
			TokenURL: "tokenEndpoint",
		},
		Scopes:      []string{"openid", "email", "profile"},
		RedirectURL: "http://localhost:8080/google/callback",
	}
	oauth := new(Oauth)
	oauth.Conn = conf
	return conf
}
