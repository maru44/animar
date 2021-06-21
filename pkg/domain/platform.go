package domain

type TPlatform struct {
	ID        int     `json:"id"`
	EngName   string  `json:"eng_name"`
	PlatName  *string `json:"plat_name"`
	BaseUrl   *string `json:"base_url"`
	Image     *string `json:"image"`
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

type TRelationPlatformInput struct {
	PlatformId int    `json:"platform_id"`
	AnimeId    int    `json:"anime_id"`
	LinkUrl    string `json:"link_url"`
}

type TPlatforms []TPlatform

type TRelationPlatforms []TRelationPlatform
