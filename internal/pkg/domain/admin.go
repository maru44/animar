package domain

type AdminInteractor interface {
	AdminAnimeInteractor
	AdminPlatformInteractor
	AdminSeasonInteractor
	AdminSeriesInteractor
	AdminCompanyInteractor
	AdminStaffInteractor
	AdminRoleInteractor
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
	PlatformId       int     `json:"platform_id"`
	AnimeId          int     `json:"anime_id"`
	LinkUrl          string  `json:"link_url"`
	DeliveryInterval *string `json:"interval,omitempty"`
	FirstBroadcast   *string `json:"first_broadcast,omitempty"`
}

type AdminPlatformInteractor interface {
	PlatformAllAdmin() (TPlatforms, error)
	PlatformDetail(int) (TPlatform, error)
	PlatformInsert(TPlatform) (int, error)
	PlatformUpdate(TPlatform, int) (int, error)
	PlatformDelete(int) (int, error)
	//relation

	RelationPlatformInsert(TRelationPlatformInput) (int, error)
	RelationPlatformUpdate(TRelationPlatformInput) (int, error)
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
	DeleteRelationSeasonAnime(animeId, seasonId int) (int, error)
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
	Name           string  `json:"name"`
	EngName        string  `json:"eng_name"`
	OfficialUrl    *string `json:"official_url,omitempty"`
	Explanation    *string `json:"explanation,omitempty"`
	TwitterAccount *string `json:"twitter_account,omitempty"`
}

type AdminCompanyInteractor interface {
	ListCompany() ([]Company, error)
	DetailCompany(string) (CompanyDetail, error)
	InsertCompany(CompanyInput) (int, error)
	UpdateCompany(CompanyInput, string) (int, error)
}

/************************
         staff
*************************/

type StaffInput struct {
	EngName    string  `json:"eng_name,omitempty"`
	FamilyName *string `json:"family_name,omitempty"`
	GivenName  *string `json:"given_name,omitempty"`
}

type AdminStaffInteractor interface {
	InsertStaff(StaffInput) (int, error)
	// UpdateStaff(StaffInput, int) (int, error)
}

/************************
         role
*************************/

type RoleInput struct {
	Num  int    `json:"num"`
	Role string `json:"role"`
}

type AnimeStaffRoleInput struct {
	AnimeId int `json:"anime_id"`
	RoleId  int `json:"role_id"`
	StaffId int `json:"staff_id"`
}

type AdminRoleInteractor interface {
	RoleList() ([]Role, error)
	InsertRole(RoleInput) (int, error)
	InsertStaffRole(AnimeStaffRoleInput) (int, error)
}
