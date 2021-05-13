package auth

import (
	"animar/v1/helper"
	"context"
	"log"

	"firebase.google.com/go/v4/auth"
)

// *auth.UserRecord
func GetUserFirebase(ctx context.Context, uid string) *auth.UserInfo {
	client := helper.FirebaseClient(ctx)
	u, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Fatalf("%s", err)
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
