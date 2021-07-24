package domain

type Company struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	EngName     string  `json:"eng_name"`
	OfficialUrl *string `json:"official_url,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
}

type CompanyDetail struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	EngName        string  `json:"eng_name"`
	OfficialUrl    *string `json:"official_url,omitempty"`
	Explanation    *string `json:"explanation,omitempty"`
	TwitterAccount *string `json:"twitter_account,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      *string `json:"updated_at,omitempty"`
}

type CompanyInteractor interface {
	ListCompany() ([]Company, error)
}
