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
