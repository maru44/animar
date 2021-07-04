package infrastructure

import (
	"animar/v1/pkg/interfaces/fires"
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type Oauth struct {
	Conn *oauth2.Config
}

type OauthToken struct {
	Token *oauth2.Token
}

type OauthClient struct {
	Client *http.Client
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
	return oauth
}

func (oauth *Oauth) Exchange(ctx context.Context, code string) (fires.OauthToken, error) {
	res := new(OauthToken)
	token, err := oauth.Conn.Exchange(ctx, code)
	res.Token = token
	return res, err
}

func (oauth *Oauth) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	res := new(OauthClient)
	client := oauth.Conn.Client(ctx, token)
	res.Client = client
	return client
}

func (token *OauthToken) Valid() bool {
	return token.Token.Valid()
}
