package domain

type TSeasonRelation struct {
	ID        int    `json:"id"`
	Year      string `json:"year"`
	Season    string `json:"season"`
	SeasonEng string `json:"season_eng"`
}

// from query params
var SeasonDict = map[string]string{
	"spring": "春",
	"summer": "夏",
	"fall":   "秋",
	"winter": "冬",
}

// from DB
var SeasonDictReverse = map[string]string{
	"春": "spring",
	"夏": "summer",
	"秋": "fall",
	"冬": "winter",
}

type SeasonInteractor interface {
	RelationSeasonByAnime(int) ([]TSeasonRelation, error)
}
