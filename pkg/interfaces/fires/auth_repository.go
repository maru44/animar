package fires

import (
	"animar/v1/configs"
	"animar/v1/pkg/domain"
	"animar/v1/pkg/tools/tools"
	"context"
	"strings"
)

type AuthRepository struct {
	Firebase
}

func (repo *AuthRepository) GetUserInfo(ctx context.Context, userId string) (uInfo domain.TUserInfo, err error) {
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

func (repo *AuthRepository) GetClaims(ctx context.Context, idToken string) (claims map[string]interface{}, err error) {
	client, err := repo.Firebase.Auth(ctx)
	token, err := client.VerifyIDToken(ctx, idToken)
	claims = token.Claims
	return
}

func (repo *AuthRepository) IsAdmin(userId string) bool {
	strAdmins := configs.AdminUsers
	admins := strings.Split(strAdmins, ", ")
	return tools.IsContainString(admins, userId)
}

func (repo *AuthRepository) GetAdminId(ctx context.Context, idToken string) string {
	claims, err := repo.GetClaims(ctx, idToken)
	if err != nil {
		return ""
	}
	id := claims["user_id"].(string)
	isAdmin := repo.IsAdmin(id)
	if !isAdmin {
		return ""
	}
	return id
}
