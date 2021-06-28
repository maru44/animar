package domain

import "context"

type TUserInfo struct {
	DisplayName string `json:"displayName,omitempty"`
	Email       string `json:"email,omitempty"`
	PhotoURL    string `json:"photoUrl,omitempty"`
	ProviderID  string `json:"providerId,omitempty"`
	UID         string `json:"rawId,omitempty"`
}

type TUserToken struct {
	Token string `json:"token,omitempty"`
}

type TLoginForm struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"` //
	Password    string `json:"password"`
	//PhotoUrl          string `json:"photoUrl"`
	ReturnSecureToken bool `json:"returnSecureToken"`
}

type TTokensForm struct {
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

type TCreateReturn struct {
	IdToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
}
type ActionCodeSettings struct {
	URL             string
	HandleCodeInApp bool
}

type TRefreshReturn struct {
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

type TProfileForm struct {
	DisplayName string `json:"displayName"`
	PhotoUrl    string `json:"photoUrl"`
}

type AuthInteractor interface {
	UserInfo(context.Context, string) (TUserInfo, error)
	Claims(context.Context, string) (map[string]interface{}, error)
	IsAdmin(string) bool
	AdminId(context.Context, string) (string, error)
	SendVerify(context.Context, string) error
	UpdateProfile(context.Context, string, TProfileForm) (TUserInfo, error)
}
