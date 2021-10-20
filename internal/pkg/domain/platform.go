package domain

type (
	TPlatform struct {
		ID        int     `json:"id"`
		EngName   string  `json:"eng_name"`
		PlatName  *string `json:"plat_name"`
		BaseUrl   *string `json:"base_url,omitempty"`
		Image     *string `json:"image,omitempty"`
		IsValid   bool    `json:"is_valid"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}

	TRelationPlatform struct {
		PlatformId       int     `json:"platform_id"`
		AnimeId          int     `json:"anime_id"`
		LinkUrl          *string `json:"link_url"`
		DeliveryInterval *string `json:"interval"`
		FirstBroadcast   *string `json:"first_broadcast"`
		CreatedAt        *string `json:"created_at"`
		UpdatedAt        *string `json:"updated_at"`
		PlatName         *string `json:"plat_name"`
	}

	RawNotificationMaterial struct {
		Platform  string
		Title     string
		Slug      string
		LinkUrl   *string
		BaseUrl   *string
		FirstTime *string
		Interval  *string
		State     string
	}

	NotificationBroadcast struct {
		Platform string  `json:"platform"` // Platform PlatName
		Title    string  `json:"title"`
		Slug     string  `json:"slug"`
		LinkUrl  *string `json:"link_url"`
		Time     *string `json:"time"`
	}

	NotifiedTargetInput struct {
		SlackID string `json:"slack_id"`
		UserID  string
	}

	TPlatforms []TPlatform

	TRelationPlatforms []TRelationPlatform

	PlatformInteractor interface {
		RelationPlatformByAnime(int) (TRelationPlatforms, error)
		// notification
		RegisterNotifiedTarget(in NotifiedTargetInput) (int, error)
		UsersChannel(userId string) (*string, error)
	}

	PlatformBatchInteractor interface {
		FilterNotificationTarget() ([]string, error)
		TargetNotificationBroadcast() ([]NotificationBroadcast, error)
		MakeSlackMessage(nbs []NotificationBroadcast) string
	}
)
