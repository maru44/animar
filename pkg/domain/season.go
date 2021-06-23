package domain

type TSeasonRelation struct {
	ID     int    `json:"id"`
	Year   string `json:"year"`
	Season string `json:"season"`
}

var SeasonDict = map[string]string{
	"spring": "春",
	"summer": "夏",
	"fall":   "秋",
	"winter": "冬",
}

type SeasonInteractor interface {
	relationSeasonByAnime(int) ([]TSeasonRelation, error)
}
