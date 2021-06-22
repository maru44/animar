package domain

type TPlatform struct {
	ID        int     `json:"id"`
	EngName   string  `json:"eng_name"`
	PlatName  *string `json:"plat_name"`
	BaseUrl   *string `json:"base_url,omitempty"`
	Image     *string `json:"image,omitempty"`
	IsValid   bool    `json:"is_valid"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type TRelationPlatform struct {
	PlatformId int     `json:"platform_id"`
	AnimeId    int     `json:"anime_id"`
	LinkUrl    *string `json:"link_url"`
	CreatedAt  *string `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
	PlatName   *string `json:"plat_name"`
}

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

type TPlatforms []TPlatform

type TRelationPlatforms []TRelationPlatform

type PlatformInteractor interface {
	RelationPlatformByAnime(int) (TRelationPlatforms, error)
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
