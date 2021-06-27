package domain

import "context"

type TUserInfo struct {
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	PhotoURL    string `json:"photoUrl,omitempty"`
	ProviderID  string `json:"providerId,omitempty"`
	UID         string `json:"rawId,omitempty"`
}

type AuthInteractor interface {
	UserInfo(context.Context, string) (TUserInfo, error)
	Claims(context.Context, string) (map[string]interface{}, error)
	IsAdmin(string) bool
	AdminId(context.Context, string) string
}
