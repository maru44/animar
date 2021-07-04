package domain

type BaseInteractor interface {
	UserId(string) (string, error)
	Claims(string) (map[string]interface{}, error)
	AdminId(string) (string, error)
	GoogleUser(string) TGoogleOauth
}
