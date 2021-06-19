package domain

type TSeasonRelation struct {
	ID     int    `json:"id"`
	Year   string `json:"year"`
	Season string `json:"season"`
}

type TSeason struct {
	ID        int    `json:"id"`
	Year      string `json:"year"`
	Season    string `json:"season"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TSeasonInput struct {
	Year   string `json:"year"`
	Season string `json:"season"`
}

var SeasonDict = map[string]string{
	"spring": "春",
	"summer": "夏",
	"fall":   "秋",
	"winter": "冬",
}
