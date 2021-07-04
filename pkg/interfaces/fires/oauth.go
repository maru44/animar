package fires

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

type Oauth interface {
	// Exchange(context.Context) (string, error)
	Client(context.Context, *oauth2.Token) *http.Client
	Exchange(context.Context, string) (OauthToken, error)
}

type OauthToken interface {
	Valid() bool
}

// type V2 interface {
// 	New(*http.Client) Tokeninfo
// }

// type TokenInfo interface {
// 	AccessToken()
// }
