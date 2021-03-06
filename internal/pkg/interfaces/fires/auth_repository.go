package fires

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/tools/mysmtp"
	"animar/v1/internal/pkg/tools/tools"
	"context"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/maru44/perr"
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
	if err != nil {
		return uInfo, perr.Wrap(err, perr.Unauthorized)
	}
	u, err := client.GetUser(ctx, userId)
	if err != nil {
		return uInfo, perr.Wrap(err, perr.Unauthorized)
	}
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
	client, err := repo.Firebase.Auth(ctx)
	if err != nil {
		return claims, perr.Wrap(err, perr.Unauthorized)
	}
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return claims, perr.Wrap(err, perr.Unauthorized)
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
		err = perr.New("", perr.Forbidden)
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
	if err != nil {
		return domain.TUserInfo{}, perr.Wrap(err, perr.Unauthorized)
	}

	var params_ auth.UserToUpdate
	if params.PhotoUrl != "" {
		params_.PhotoURL(params.PhotoUrl)
	}
	params_.DisplayName(params.DisplayName)
	u, err := client.UpdateUser(ctx, userId, &params_)
	if err != nil {
		return domain.TUserInfo{}, perr.Wrap(err, perr.BadRequest)
	}

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
		return userId, perr.Wrap(err, perr.BadRequest)
	}
	userId = claims["user_id"].(string)
	return
}
