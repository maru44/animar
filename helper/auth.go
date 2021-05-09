package helper

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func FirebaseApp(ctx context.Context) *firebase.App {
	opt := option.WithCredentialsFile("secret_key.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Panic(err.Error())
	}
	return app
}

func FirebaseClient(ctx context.Context) *auth.Client {
	app := FirebaseApp(ctx)
	client, err := app.Auth(ctx)
	if err != nil {
		log.Panic(err.Error())
	}
	return client
}
