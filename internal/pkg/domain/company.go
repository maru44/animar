package domain

type Company struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	EngName     string  `json:"eng_name"`
	OfficialUrl *string `json:"official_url,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
}

type CompanyInput struct {
	Name        string  `json:"name"`
	EngName     string  `json:"eng_name"`
	OfficialUrl *string `json:"official_url,omitempty"`
}

type CompanyInteractor interface {
	Insert(CompanyInput) (int, error)
}
