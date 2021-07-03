package fires

import (
	"context"

	"firebase.google.com/go/v4/auth"
)

type Firebase interface {
	Auth(context.Context) (Client, error)
	//FirebaseAuth
}

type Client interface {
	// VerifyIDToken(context.Context, string) (Token, error)
	VerifyIDToken(context.Context, string) (*auth.Token, error)
	GetUser(context.Context, string) (*auth.UserRecord, error)
	EmailVerificationLinkWithSettings(context.Context, string, *auth.ActionCodeSettings) (string, error)
	UpdateUser(context.Context, string, *auth.UserToUpdate) (*auth.UserRecord, error)
}

type Token struct {
	Claims map[string]interface{}
}

type UserRecode interface {
	UserInfo
}

type UserInfo interface{}

type ActionCodeSettings interface{}

type UserToUpdate interface{}
