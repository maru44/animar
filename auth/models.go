package auth

import (
	"animar/v1/helper"
	"context"
	"fmt"

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
		//log.Fatalf("%s", err)
		fmt.Print(err)
		return nil
	}
	claims := token.Claims

	return claims
}

// 不要では?
func CreateUser(
	ctx context.Context, client *auth.Client,
	email string, password string, // name string,
) *auth.UserRecord {
	params := (&auth.UserToCreate{}).
		Email(email).
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

	// @TODO Use env!!
	actionCodeSettings := &auth.ActionCodeSettings{
		URL:             "http://localhost:3000/auth/verify/",
		HandleCodeInApp: false,
	}

	link, err := client.EmailVerificationLinkWithSettings(ctx, email, actionCodeSettings)
	fmt.Print(link)

	return u
}

// send email link
func SendVerifyEmailAtRegister(ctx context.Context, client *auth.Client, email string) error {
	// @TODO Use env!!
	actionCodeSettings := &auth.ActionCodeSettings{
		URL:             "http://localhost:3000/",
		HandleCodeInApp: false,
	}

	link, err := client.EmailVerificationLinkWithSettings(ctx, email, actionCodeSettings)
	if err != nil {
		fmt.Print(err.Error())
	}

	sended := helper.SendVerifyEmail(email, link)
	return sended
}

func SetAdminClaim(ctx context.Context, client *auth.Client, uid string) {
	claims := map[string]interface{}{"is_admin": false}
	err := client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		panic(err.Error())
	}
}

// 不要
func VerifyEmail(ctx context.Context, client *auth.Client, uid string) {
	params := (&auth.UserToUpdate{}).
		EmailVerified(true)
	u, err := client.UpdateUser(ctx, uid, params)
	if err != nil {
		panic(err.Error())
	}
	fmt.Print(u)
}
