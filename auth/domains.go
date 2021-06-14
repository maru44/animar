package auth

import (
	"animar/v1/configs"
	"animar/v1/tools/fire"
	"animar/v1/tools/mysmtp"
	"animar/v1/tools/tools"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"firebase.google.com/go/v4/auth"
)

func GetJWTFromGoogle(posted TLoginForm) TTokensForm {
	var tokens TTokensForm

	posted_json, _ := json.Marshal(posted)
	url := `https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=` + configs.FirebaseApiKey
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(posted_json),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tools.ErrorLog(err)
		return tokens
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &tokens)

	if err != nil {
		tools.ErrorLog(err)
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
	protocol := "http://"
	if tools.IsProductionEnv() {
		protocol = "https://"
	}
	actionCodeSettings := &auth.ActionCodeSettings{
		URL:             protocol + configs.FrontHost + configs.FrontPort + "/auth" + "/confirmed",
		HandleCodeInApp: false,
	}

	link, err := client.EmailVerificationLinkWithSettings(ctx, email, actionCodeSettings)
	if err != nil {
		fmt.Print(err.Error())
	}

	sended := mysmtp.SendVerifyEmail(email, link)
	return sended
}

func SetAdminClaimExe(ctx context.Context, client *auth.Client, idToken *string) {
	userId := fire.GetUserIdFromToken(*idToken)
	SetAdminClaim(ctx, client, userId)
}
