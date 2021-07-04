package fires

import (
	"context"

	"golang.org/x/oauth2"
)

type Oauth interface {
	Exchange(context.Context) (string, error)
	Client(context.Context, *oauth2.Token)
}
