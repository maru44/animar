package domain

type TAnime struct {
	ID            int     `json:"id"`
	Slug          string  `json:"slug"`
	Title         string  `json:"title"`
	ThumbUrl      *string `json:"thumb_url,omitempty"`
	CopyRight     *string `json:"copyright,omitempty"`
	Abbreviation  *string `json:"abbreviation,omitempty"`
	Description   *string `json:"description,omitempty"`
	State         *string `json:"state,omitempty"`
	SeriesId      *int    `json:"series_id,omitempty"`
	CountEpisodes *int    `json:"count_episodes,omitemptys"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
}

type TAnimeMinimum struct {
	ID    int    `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type TAnimeWithSeries struct {
	ID            int     `json:"id"`
	Slug          string  `json:"slug"`
	Title         string  `json:"title"`
	ThumbUrl      *string `json:"thumb_url,omitempty"`
	CopyRight     *string `json:"copyright,omitempty"`
	Abbreviation  *string `json:"abbreviation,omitempty"`
	Description   *string `json:"description,omitempty"`
	State         *string `json:"state,omitempty"`
	SeriesId      *int    `json:"series_id,omitempty"`
	CountEpisodes *int    `json:"count_episodes,omitemptys"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
	SeriesName    *string `json:"series_name"`
}

type TAnimes []TAnime

type TAnimeMinimums []TAnimeMinimum

type TAnimeAdmin struct {
	ID            int     `json:"id"`
	Slug          string  `json:"slug"`
	Title         string  `json:"title"`
	Abbreviation  *string `json:"abbreviation,omitempty"`
	Kana          *string `json:"kana,omitempty"`
	EngName       *string `json:"eng_name,omitempty"`
	ThumbUrl      *string `json:"thumb_url,omitempty"`
	CopyRight     *string `json:"copyright,omitempty"`
	Description   *string `json:"description,omitempty"`
	State         *string `json:"state,omitempty"`
	SeriesId      *int    `json:"series_id,omitempty"`
	CountEpisodes *int    `json:"count_episodes,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
}

type TAnimeAdmins []TAnimeAdmin

// type TAnimeInput struct {
// 	Title         string `json:"title"`
// 	Abbrevation   string `json:"abbrevation,omitempty"`
// 	Kana          string `json:"kana,omitempty"`
// 	EngName       string `json:"eng_name:omitempty"`
// 	ThumbUrl      string `json:"thumb_url,omitempty"`
// 	PreThumbUrl   string `json:"pre_thumb,omitempty"`
// 	Description   string `json:"description,omitempty"`
// 	State         int    `json:"state,omitempty"`
// 	SeriesId      int    `json:"series_id,omitempty"`
// 	CountEpisodes int    `jsoin:"count_episodes,omitempty"`
// }

type TAnimeInsert struct {
	Title         string  `json:"title"`
	Slug          string  `json:"slug"`
	Abbreviation  *string `json:"abbrevation,omitempty"`
	Kana          *string `json:"kana,omitempty"`
	EngName       *string `json:"eng_name:omitempty"`
	Description   *string `json:"description,omitempty"`
	ThumbUrl      *string `json:"thumb_url,omitempty"`
	State         *string `json:"state,omitempty"`
	SeriesId      *int    `json:"series_id,omitempty"`
	CountEpisodes *int    `json:"count_episodes,omitempty"`
	Copyright     *string `json:"copyright,omitempty"`
}

func (a *TAnime) GetId() int {
	return a.ID
}

func (a *TAnimeWithSeries) GetId() int {
	return a.ID
}
