package infrastructure

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/fires"
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Firebase struct {
	App *firebase.App
}

type FirebaseClient struct {
	Client *auth.Client
}

type FirebaseToken struct {
	Token *auth.Token
}

type FirebaseRecord struct {
	UserRecord *auth.UserRecord
}

type FirebaseUserInfo struct {
	UserInfo *auth.UserInfo
}

type FirebaseSettings struct {
	ActionCodeSettings *auth.ActionCodeSettings
}

type FirebaseUpdate struct {
	UserToUpdate *auth.UserToUpdate
}

// initialize firebase sdk
func NewFireBaseClient() fires.Firebase {
	ctx := context.Background()
	opt := option.WithCredentialsFile("../../configs/secret_key.json")
	config := &firebase.Config{ProjectID: configs.ProjectId} // for Google Oauth
	app, _ := firebase.NewApp(ctx, config, opt)
	firebase_ := new(Firebase)
	firebase_.App = app
	return firebase_
}

func (f *Firebase) Auth(ctx context.Context) (fires.Client, error) {
	client := new(FirebaseClient)
	res, err := f.App.Auth(ctx)
	if err != nil {
		return client, domain.Errors{Inner: err, Flag: domain.FirebaseConnectionError}
	}
	client.Client = res
	return client, nil
}

// あと少し *auth.Tokenと型が違うからだめっぽ *auth.Tokenに合わせてfieldを追加してみたがダメ     (fires.Token, error)
func (fc *FirebaseClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token := new(FirebaseToken)
	realToken, err := fc.Client.VerifyIDToken(ctx, idToken)
	if err != nil {
		domain.ErrorWarn(err)
	}
	token.Token = realToken
	return token.Token, err
}

func (fc *FirebaseClient) GetUser(ctx context.Context, userId string) (*auth.UserRecord, error) {
	return fc.Client.GetUser(ctx, userId)
}

func (fc *FirebaseClient) EmailVerificationLinkWithSettings(ctx context.Context, email string, settings *auth.ActionCodeSettings) (string, error) {
	return fc.Client.EmailVerificationLinkWithSettings(ctx, email, settings)
}

func (fc *FirebaseClient) UpdateUser(ctx context.Context, userId string, user *auth.UserToUpdate) (*auth.UserRecord, error) {
	return fc.Client.UpdateUser(ctx, userId, user)
}

// 不要では?
func (token *FirebaseToken) Claims() map[string]interface{} {
	return token.Token.Claims
}
