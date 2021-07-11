package domain

type TSeries struct {
	ID         int     `json:"id"`
	EngName    string  `json:"eng_name"`
	SeriesName *string `json:"series_name,omitempty"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  *string `json:"updated_at,omitempty"`
}
