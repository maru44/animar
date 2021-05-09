package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
)

func UserListView(w http.ResponseWriter, r *http.Request) {
	jsonStr := `{"grant_type": "client_credentials", "client_id": "` + os.Getenv("AUTH0_CLIENT_ID") + `", "client_secret": "` + os.Getenv("AUTH0_SECRET") + `", "audience": "https://` + os.Getenv("AUTH0_DOMAIN") + `/api/v2/"}`
	url := `https://` + os.Getenv("AUTH0_DOMAIN") + `/oauth/token`
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(jsonStr)),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	fmt.Fprintln(w, resp)
}

func SampleGetUser(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uid := query.Get("uid")
	ctx := context.Background()

	user := GetUserFirebase(ctx, uid)
	fmt.Fprintln(w, user.Email, user.EmailVerified, user.UID, user.DisplayName)
}
