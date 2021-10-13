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
	PlatformId       int     `json:"platform_id"`
	AnimeId          int     `json:"anime_id"`
	LinkUrl          *string `json:"link_url"`
	DeliveryInterval *string `json:"interval"`
	FirstBroadcast   *string `json:"first_broadcast"`
	CreatedAt        *string `json:"created_at"`
	UpdatedAt        *string `json:"updated_at"`
	PlatName         *string `json:"plat_name"`
}

type RawNotificationMaterial struct {
	Platform  string
	Title     string
	LinkUrl   *string
	BaseUrl   *string
	FirstTime *string
	Interval  *string
	State     string
}

type NotificationBroadcast struct {
	Platform string  `json:"platform"` // Platform PlatName
	Title    string  `json:"title"`
	LinkUrl  *string `json:"link_url"`
	Time     *string `json:"time"`
}

type TPlatforms []TPlatform

type TRelationPlatforms []TRelationPlatform

type PlatformInteractor interface {
	RelationPlatformByAnime(int) (TRelationPlatforms, error)
	TargetNotificationBroadcast() ([]NotificationBroadcast, error)
}
