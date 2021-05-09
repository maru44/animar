package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	jwt "github.com/dgrijalva/jwt-go"
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

func GetJWTPayload(token string) {
	//hCS := strings.split(token, ".")
	parsedToken, err := jwt.Parse(token)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print(parsedToken)
}
