package tools

import (
	"animar/v1/configs"
	"context"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type TUserIdCookieInput struct {
	Token string `json:"token"`
}

func FirebaseApp(ctx context.Context) *firebase.App {
	opt := option.WithCredentialsFile("../../configs/secret_key.json")
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

// userID from idToken
func GetUserIdFromToken(idToken string) string {
	ctx := context.Background()
	client := FirebaseClient(ctx)
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

// post や update deleteは
// これがblankであれば問答無用で400
func GetIdFromCookie(r *http.Request) string {
	idToken, err := r.Cookie("idToken")
	if err != nil {
		//fmt.Print(err.Error())
		return ""
	}
	id := GetUserIdFromToken(idToken.Value)

	return id
}

func IsAdmin(userId string) bool {
	strAdmins := configs.AdminUsers
	admins := strings.Split(strAdmins, ", ")

	return IsContainString(admins, userId)
}

// isAdmin required
// anonymous & not admin >> blank
// else >> return userId
func GetAdminIdFromCookie(r *http.Request) string {
	idToken, err := r.Cookie("idToken")
	if err != nil {
		return ""
	}
	id := GetUserIdFromToken(idToken.Value)
	isAdmin := IsAdmin(id)
	if !isAdmin {
		return ""
	}
	return id
}

// is admin
func GetAdminIdFromIdToken(idToken string) string {
	id := GetUserIdFromToken(idToken)
	isAdmin := IsAdmin(id)
	if !isAdmin {
		return ""
	}
	return id
}
