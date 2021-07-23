package domain

type AdminInteractor interface {
	AdminAnimeInteractor
	AdminPlatformInteractor
	AdminSeasonInteractor
	AdminSeriesInteractor
	AdminCompanyInteractor
}

/************************
         anime
*************************/

type AnimeAdmin struct {
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
	CompanyId     *int    `json:"company_id,omitempty"`
	HashTag       *string `json:"hash_tag,omitempty"`
	TwitterUrl    *string `json:"twitter_url,omitempty"`
	OfficialUrl   *string `json:"official_url,omitempty"`
}

type AnimeAdmins []AnimeAdmin

type AnimeInsert struct {
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
	CompanyId     *int    `json:"company_id,omitempty"`
	HashTag       *string `json:"hash_tag,omitempty"`
	TwitterUrl    *string `json:"twitter_url,omitempty"`
	OfficialUrl   *string `json:"official_url,omitempty"`
}

type AdminAnimeInteractor interface {
	AnimesAllAdmin() (TAnimes, error)
	AnimeDetailAdmin(int) (AnimeAdmin, error)
	AnimeInsert(AnimeInsert) (int, error)
	AnimeUpdate(int, AnimeInsert) (int, error)
	AnimeDelete(int) (int, error)
}

/************************
         platform
*************************/

type TPlatformInput struct {
	EngName  string  `json:"eng_name"`
	PlatName *string `json:"plat_name,omitempty"`
	BaseUrl  *string `json:"base_url,omitempty"`
	Image    *string `json:"image,omitempty"`
	IsValid  bool    `json:"is_valid"`
}

type TRelationPlatformInput struct {
	PlatformId int    `json:"platform_id"`
	AnimeId    int    `json:"anime_id"`
	LinkUrl    string `json:"link_url"`
}

type AdminPlatformInteractor interface {
	PlatformAllAdmin() (TPlatforms, error)
	PlatformDetail(int) (TPlatform, error)
	PlatformInsert(TPlatform) (int, error)
	PlatformUpdate(TPlatform, int) (int, error)
	PlatformDelete(int) (int, error)
	RelationPlatformInsert(TRelationPlatformInput) (int, error)
	RelationPlatformDelete(int, int) (int, error)
	RelationPlatformByAnime(int) (TRelationPlatforms, error)
}

/************************
         season
*************************/

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

type TSeasonRelationInput struct {
	SeasonId int `json:"season_id"`
	AnimeId  int `json:"anime_id"`
}

type AdminSeasonInteractor interface {
	InsertSeason(TSeasonInput) (int, error)
	ListSeason() ([]TSeason, error)
	DetailSeason(int) (TSeason, error)
	InsertRelationSeasonAnime(TSeasonRelationInput) (int, error)
}

/************************
         series
*************************/

type TSeriesInput struct {
	EngName    string `json:"eng_name"`
	SeriesName string `json:"series_name,omitempty"`
}

type AdminSeriesInteractor interface {
	ListSeries() ([]TSeries, error)
	DetailSeries(int) (TSeries, error)
	InsertSeries(TSeriesInput) (int, error)
	UpdateSeries(TSeriesInput, int) (int, error)
	DeleteSeries(int) (int, error)
}

/************************
         company
*************************/

type CompanyInput struct {
	Name        string  `json:"name"`
	EngName     string  `json:"eng_name"`
	OfficialUrl *string `json:"official_url,omitempty"`
}

type AdminCompanyInteractor interface {
	InsertCompany(CompanyInput) (int, error)
}
