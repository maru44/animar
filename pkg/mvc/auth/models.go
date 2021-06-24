package auth

import (
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/tools/tools"
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

// *auth.UserRecord
func GetUserFirebase(ctx context.Context, uid string) *auth.UserInfo {
	client := fire.FirebaseClient(ctx)
	u, err := client.GetUser(ctx, uid)
	if err != nil {
		//fmt.Println(err.Error(), err)
		return nil
	}
	return u.UserInfo
}

// user info from token
func VerifyFirebase(ctx context.Context, idToken string) map[string]interface{} {
	client := fire.FirebaseClient(ctx)
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		tools.ErrorLog(err)
		return nil
	}
	claims := token.Claims

	return claims
}

func SetAdminClaim(ctx context.Context, client *auth.Client, uid string) {
	claims := map[string]interface{}{"is_admin": false}
	err := client.SetCustomUserClaims(ctx, uid, claims)
	if err != nil {
		fmt.Print("failure set custom claims")
	}
}

//
/*
func VerifyEmail(ctx context.Context, client *auth.Client, uid string) {
	params := (&auth.UserToUpdate{}).
		EmailVerified(true)
	u, err := client.UpdateUser(ctx, uid, params)
	if err != nil {
		panic(err.Error())
	}
	fmt.Print(u)
}
*/
