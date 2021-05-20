package helper

import (
	"context"
	"log"
	"net/http"
	"strings"

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

func GetClaimsFromCookie(r *http.Request) map[string]interface{} {
	idToken, err := r.Cookie("idToken")
	if err != nil {
		return nil
	}

	ctx := context.Background()
	client := FirebaseClient(ctx)
	token, err := client.VerifyIDToken(ctx, idToken.Value)
	if err != nil {
		//log.Fatalf("%s", err)
		return nil
	}
	claims := token.Claims

	return claims
}

func GetIdFromCookie(r *http.Request) string {
	idToken, err := r.Cookie("idToken")
	if err != nil {
		//fmt.Print(err.Error())
		return ""
	}

	ctx := context.Background()
	client := FirebaseClient(ctx)
	// @TODO getUserIdFromToken使う
	token, err := client.VerifyIDToken(ctx, idToken.Value)
	if err != nil {
		//fmt.Printf("%s%s", err.Error(), err)
		if strings.Contains(err.Error(), "ID token has expired at:") {
			return ""
		}
	}
	claims := token.Claims
	id := claims["user_id"]

	return id.(string)
}

func GetUserIdFromToken(ctx context.Context, client *auth.Client, idToken string) string {
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		//fmt.Printf("%s%s", err.Error(), err)
		if strings.Contains(err.Error(), "ID token has expired at:") {
			return ""
		}
	}
	claims := token.Claims
	id := claims["user_id"]

	return id.(string)
}
