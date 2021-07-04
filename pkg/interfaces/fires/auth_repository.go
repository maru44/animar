package fires

import (
	"animar/v1/configs"
	"animar/v1/pkg/domain"
	"animar/v1/pkg/tools/mysmtp"
	"animar/v1/pkg/tools/tools"
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"firebase.google.com/go/v4/auth"
	"golang.org/x/oauth2"
)

type AuthRepository struct {
	Firebase
}

const (
	authorizeEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenEndpoint     = "https://www.googleapis.com/oauth2/v4/token"
)

func (repo *AuthRepository) GetUserInfo(userId string) (uInfo domain.TUserInfo, err error) {
	ctx := context.Background()
	client, err := repo.Firebase.Auth(ctx)
	u, err := client.GetUser(ctx, userId)
	info := u.UserInfo
	uInfo = domain.TUserInfo{
		DisplayName: info.DisplayName,
		Email:       info.Email,
		PhotoURL:    info.PhotoURL,
		ProviderID:  info.ProviderID,
		UID:         info.UID,
	}
	return
}

func (repo *AuthRepository) GetClaims(idToken string) (claims map[string]interface{}, err error) {
	ctx := context.Background()
	client, _ := repo.Firebase.Auth(ctx)
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return
	}
	claims = token.Claims
	return
}

func (repo *AuthRepository) IsAdmin(userId string) bool {
	strAdmins := configs.AdminUsers
	admins := strings.Split(strAdmins, ", ")
	return tools.IsContainString(admins, userId)
}

func (repo *AuthRepository) GetAdminId(idToken string) (userId string, err error) {
	claims, err := repo.GetClaims(idToken)
	if err != nil {
		return
	}
	userId = claims["user_id"].(string)
	isAdmin := repo.IsAdmin(userId)
	if !isAdmin {
		err = domain.ErrForbidden
		return
	}
	return
}

func (repo *AuthRepository) SendVerifyEmail(email string) (err error) {
	protocol := "http://"
	if tools.IsProductionEnv() {
		protocol = "https://"
	}

	ctx := context.Background()
	client, err := repo.Firebase.Auth(ctx)
	settings := &auth.ActionCodeSettings{
		URL:             protocol + configs.FrontHost + configs.FrontPort + "/auth" + "/confirmed",
		HandleCodeInApp: false,
	}
	link, err := client.EmailVerificationLinkWithSettings(ctx, email, settings)
	err = mysmtp.SendVerifyEmail(email, link)
	return err
}

func (repo *AuthRepository) Update(userId string, params domain.TProfileForm) (domain.TUserInfo, error) {
	ctx := context.Background()
	client, err := repo.Firebase.Auth(ctx)
	var params_ auth.UserToUpdate
	if params.PhotoUrl != "" {
		params_.DisplayName(params.DisplayName)
	}
	params_.PhotoURL(params.DisplayName)
	u, err := client.UpdateUser(ctx, userId, &params_)
	user := domain.TUserInfo{
		DisplayName: u.UserInfo.DisplayName,
		Email:       u.UserInfo.Email,
		PhotoURL:    u.UserInfo.PhotoURL,
		ProviderID:  u.UserInfo.ProviderID,
		UID:         u.UserInfo.UID,
	}
	return user, err
}

func (repo *AuthRepository) GetUserId(idToken string) (userId string, err error) {
	claims, err := repo.GetClaims(idToken)
	if err != nil {
		return
	}
	userId = claims["user_id"].(string)
	return
}

// 仮のベタ書き
func (repo *AuthRepository) GoogleOAuth() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     configs.GoogleWebClientId,
		ClientSecret: configs.GoogleWebClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizeEndpoint,
			TokenURL: tokenEndpoint,
		},
		Scopes: []string{
			// "https://www.googleapis.com/auth/firebase",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			//"https://www.googleapis.com/auth/cloud-platform", // 全然ダメ
		},
		RedirectURL: "http://localhost:8000/auth/google/redirect",
	}
	return config
}

func (repo *AuthRepository) GoogleOAuthCallback() {
	config := repo.GoogleOAuth()
	// ctx := context.Background()
	// var tok *oauth2.Token
	// tok, err := config.Exchange(ctx, "aaa")
	// fmt.Print(tok)
	// if err != nil {
	// 	tools.ErrorLog(err)
	// }
	// if tok.Valid() == false {
	// }
	fmt.Print(config.AuthCodeURL("aaa"))

	// service, _ := v2.New(config.Client(ctx, tok))
	//var tokenInfo *v2.Tokeninfo
	// tokenInfo, _ = service.Tokeninfo().AccessToken(tok.AccessToken).Context(ctx).Do()
	// fmt.Print(tokenInfo)
}

func (repo *AuthRepository) GoogleRedirect(code string) {
	config := repo.GoogleOAuth()
	ctx := context.Background()
	var tok *oauth2.Token
	tok, err := config.Exchange(ctx, code)

	if err != nil {
		tools.ErrorLog(err)
	}
	client := config.Client(ctx, tok)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tok.AccessToken)
	if err != nil {
		tools.ErrorLog(err)
	}
	defer res.Body.Close()
	byteArray, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(byteArray))
}
