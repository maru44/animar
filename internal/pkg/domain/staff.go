package domain

type Staff struct {
	ID         int     `json:"id"`
	EngName    string  `json:"eng_name"`
	FamilyName *string `json:"family_name,omitempty"`
	GivenName  *string `json:"given_name,omitempty"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type StaffInteractor interface {
	StaffList() ([]Staff, error)
	// StaffDetail() (Staff, error)
}
