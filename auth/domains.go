package auth

import (
	"animar/v1/tools"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"firebase.google.com/go/v4/auth"
)

func GetJWTFromGoogle(posted TLoginForm) TTokensForm {
	var tokens TTokensForm

	posted_json, _ := json.Marshal(posted)
	url := `https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=` + os.Getenv("FIREBASE_API_KEY")
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(posted_json),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		return tokens
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &tokens)

	if err != nil {
		fmt.Print(err.Error())
		return tokens
	}

	return tokens
}

/*
func GetJWTPayload(token string) {
	//hCS := strings.split(token, ".")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print(parsedToken)
}
*/

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

	sended := tools.SendVerifyEmail(email, link)
	return sended
}

func SetAdminClaimExe(ctx context.Context, client *auth.Client, idToken *string) {
	userId := tools.GetUserIdFromToken(*idToken)
	fmt.Print(userId)
	SetAdminClaim(ctx, client, userId)
}
