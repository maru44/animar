package queryset

const (
	NotificationTarget = "SELECT GROUP_CONCAT(slack_channel) " +
		"FROM slack_channel " +
		"WHERE is_active = true"

	TodaysBroadcastQuery = "SELECT plat.plat_name, animes.title, animes.slug, rel.link_url, plat.base_url, rel.first_broadcast, rel.delivery_interval, animes.state " +
		"FROM relation_anime_platform AS rel " +
		"LEFT JOIN animes ON rel.anime_id = animes.id " +
		"LEFT JOIN platforms AS plat ON rel.platform_id = plat.id " +
		"WHERE animes.state = 'now' " +
		"OR ((animes.state = 'pre' OR (rel.delivery_interval = 'once' AND animes.state NOT IN ('now','pre'))) AND " +
		"rel.first_broadcast BETWEEN DATE_ADD(DATE(NOW()), INTERVAL 4 HOUR) " +
		"AND DATE_ADD(DATE(NOW()), INTERVAL 28 HOUR))"

	PlatformFilterByAnimeQuery = "SELECT relation_anime_platform.*, platforms.plat_name FROM relation_anime_platform " +
		"LEFT JOIN platforms ON relation_anime_platform.platform_id = platforms.id " +
		"WHERE anime_id = ?"

	RegisterNotifiedTargetQuery = "INSERT INTO slacks(slack_channel, user_id) VALUES (?, ?)"

	GetUsersSlackChannelQuery = "SELECT slack_channel FROM slacks WHERE user_id = ?"
)
