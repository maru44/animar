package infrastructure

import (
	"animar/v1/configs"
	"animar/v1/pkg/interfaces/fires"
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
	client, err := f.App.Auth(ctx)
	if err != nil {
		return new(FirebaseClient), err
	}
	return client, nil
}

func (fc *FirebaseClient) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return fc.Client.VerifyIDToken(ctx, idToken)
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
