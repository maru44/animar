package auth

import (
	"animar/v1/helper"
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/v4/auth"
)

// *auth.UserRecord
func GetUserFirebase(ctx context.Context, uid string) *auth.UserInfo {
	client := helper.FirebaseClient(ctx)
	u, err := client.GetUser(ctx, uid)
	if err != nil {
		fmt.Println(err.Error(), err)
	}
	return u.UserInfo
}

// user info from token
func VerifyFirebase(ctx context.Context, idToken string) map[string]interface{} {
	client := helper.FirebaseClient(ctx)
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("%s", err)
	}
	claims := token.Claims

	return claims
}

// 不要では?
func CreateUser(
	ctx context.Context, client *auth.Client,
	emai string, password string, // name string,
) *auth.UserRecord {
	params := (&auth.UserToCreate{}).
		Email(emai).
		EmailVerified(false).
		Password(password).
		//PhoneNumber(phoneNumber).
		//DisplayName(name).
		//PhotoURL(photo).
		Disabled(false)
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		panic(err.Error())
	}

	return u
}

func SetAdminClaim(ctx context.Context, client *auth.Client, uid string) {
	claims := map[string]interface{}{"is_admin": false}
	err := client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		panic(err.Error())
	}
}

func VerifyEmail(ctx context.Context, client *auth.Client, uid string) {
	params := (&auth.UserToUpdate{}).
		EmailVerified(true)
	u, err := client.UpdateUser(ctx, uid, params)
	if err != nil {
		panic(err.Error())
	}
	fmt.Print(u)
}
