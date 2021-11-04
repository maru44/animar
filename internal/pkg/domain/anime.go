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
	// detail of anime
	ID             int     `json:"id"`
	Slug           string  `json:"slug"`
	Title          string  `json:"title"`
	ThumbUrl       *string `json:"thumb_url,omitempty"`
	CopyRight      *string `json:"copyright,omitempty"`
	Abbreviation   *string `json:"abbreviation,omitempty"`
	Description    *string `json:"description,omitempty"`
	State          *string `json:"state,omitempty"`
	SeriesId       *int    `json:"series_id,omitempty"`
	CompanyId      *int    `json:"company_id,omitempty"`
	CountEpisodes  *int    `json:"count_episodes,omitemptys"`
	HashTag        *string `json:"hash_tag,omitempty"`
	TwitterUrl     *string `json:"twitter_url,omitempty"`
	OfficialUrl    *string `json:"official_url,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      *string `json:"updated_at,omitempty"`
	SeriesName     *string `json:"series_name,omitempty"`
	CompanyName    *string `json:"company_name,omitempty"`
	CompanyEngName *string `json:"company_eng_name,omitempty"`
}

type TAnimes []TAnime

type TAnimeMinimums []TAnimeMinimum

func (a *TAnime) GetId() int {
	return a.ID
}

func (a *TAnimeWithSeries) GetId() int {
	return a.ID
}

type AnimeInteractor interface {
	AnimesAll() (TAnimes, error)
	AnimesOnAir() (TAnimes, error)
	AnimesSearch(string) (TAnimes, error)
	AnimesBySeason(string, string) (TAnimes, error)
	AnimesBySeries(id int) ([]TAnimeWithSeries, error)
	AnimesByCompany(string) (TAnimes, error)
	AnimeMinimums() (TAnimeMinimums, error)
	AnimeSearchMinimum(string) (TAnimeMinimums, error)
	/*   detail   */
	AnimeDetail(int) (TAnime, error)
	AnimeDetailBySlug(string) (TAnimeWithSeries, error)
	/*   review   */
	ReviewFilterByAnime(int, string) (TReviews, error)
	/*   company   */
	DetailCompanyByEng(string) (CompanyDetail, error)
	/*   change on air   */
}

const (
	StateNow = "now"
	StatePre = "pre"
	StateCut = "cut"
	StateFin = "fin"
)
