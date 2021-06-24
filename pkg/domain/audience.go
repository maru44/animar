package domain

type TAudience struct {
	ID        int    `json:"id"`
	State     int    `json:"state"`
	AnimeId   int    `json:"anime_id"`
	UserId    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TAudienceCount struct {
	State int `json:"state"`
	Count int `json:"count"`
}

type TAudienceJoinAnime struct {
	ID        int     `json:"id"`
	State     int     `json:"state"`
	AnimeId   int     `json:"anime_id"`
	UserId    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Title     string  `json:"title"`
	Slug      string  `json:"slug"`
	Content   *string `json:"content"`
	AState    *string `json:"anime_state"`
}

type TAudienceInput struct {
	AnimeId int `json:"anime_id"`
	State   int `json:"state,string"` // form
}

type AudienceInteractor interface {
	AnimeAudienceCounts(int) ([]TAudienceCount, error)
	//AudienceByUser(string) ([]TAudience, error)
	AudienceWithAnimeByUser(string) ([]TAudienceJoinAnime, error)
	AudienceByAnimeAndUser(int, string) (TAudience, error)
	InsertAudience(TAudienceInput, string) (int, error)
	UpsertAudience(TAudienceInput, string) (int, error)
	DeleteAudience(int, string) (int, error)
}
