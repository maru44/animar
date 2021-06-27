package fires

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

// type Option interface {
// 	WithCredentialsFile(string) ClientOption
// }

// type ClientOption interface {
// 	Apply(DialSettings)
// }

// app

type Firebase interface {
	Auth(context.Context) (Client, error)
}

type Client interface {
	VerifyIDToken(context.Context, string) (*auth.Token, error)
	GetUser(context.Context, string) (*auth.UserRecord, error)
	EmailVerificationLinkWithSettings(context.Context, string, *auth.ActionCodeSettings) (string, error)
	UpdateUser(context.Context, string, *auth.UserToUpdate) (*auth.UserRecord, error)
}
