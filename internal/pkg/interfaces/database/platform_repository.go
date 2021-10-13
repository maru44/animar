package database

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/tools/tools"
	"strings"
	"time"

	"github.com/maru44/perr"
)

type PlatformRepository struct {
	SqlHandler
}

func (repo *PlatformRepository) FilterByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	rows, err := repo.Query(
		"Select relation_anime_platform.*, platforms.plat_name FROM relation_anime_platform "+
			"LEFT JOIN platforms ON relation_anime_platform.platform_id = platforms.id "+
			"WHERE anime_id = ?", animeId,
	)
	if err != nil {
		return platforms, perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.TRelationPlatform
		err = rows.Scan(
			&p.PlatformId, &p.AnimeId, &p.LinkUrl,
			&p.DeliveryInterval, &p.FirstBroadcast,
			&p.CreatedAt, &p.UpdatedAt, &p.PlatName,
		)
		if err != nil {
			return platforms, perr.Wrap(err, perr.NotFound)
		}
		platforms = append(platforms, p)
	}
	return
}

func (repo *PlatformRepository) FilterTodaysBroadCast() ([]domain.NotificationBroadcast, error) {
	rows, err := repo.Query(
		"SELECT plat.plat_name, animes.title, rel.link_url, plat.base_url, rel.first_broadcast, rel.delivery_interval, animes.state " +
			"FROM relation_anime_platform AS rel " +
			"LEFT JOIN animes ON rel.anime_id = animes.id " +
			"LEFT JOIN platforms AS plat ON rel.platform_id = plat.id " +
			"WHERE animes.state = 'now' " +
			"OR ((animes.state = 'pre' OR (rel.delivery_interval = 'once' AND animes.state NOT IN ('now','pre'))) AND " +
			"rel.first_broadcast BETWEEN DATE_ADD(DATE(NOW()), INTERVAL 4 HOUR)) " +
			"AND DATE_ADD(DATE(NOW()), INTERVAL 28 HOUR))",
	)
	if err != nil {
		return nil, perr.Wrap(err, perr.BadRequest)
	}

	if !rows.Next() {
		return nil, perr.New("", perr.NotFound)
	}

	var out []domain.NotificationBroadcast
	for rows.Next() {
		var m domain.RawNotificationMaterial
		err := rows.Scan(
			&m.Platform, &m.Title, &m.LinkUrl, &m.BaseUrl,
			&m.FirstTime, &m.Interval, &m.State,
		)
		if err != nil {
			return nil, perr.Wrap(err, perr.BadRequest)
		}

		b := repo.modifyBroadcastTime(&m)
		out = append(out, *b)
	}

	return out, nil
}

func (repo *PlatformRepository) modifyBroadcastTime(m *domain.RawNotificationMaterial) *domain.NotificationBroadcast {
	if m.Interval == nil {
		return nil
	}

	switch m.State {
	case "now":
		if *m.Interval == "once" {
			return &domain.NotificationBroadcast{
				Platform: m.Platform,
				Title:    m.Title,
				LinkUrl:  m.LinkUrl,
				Time:     extractFristTime(m.FirstTime),
			}
		} else if *m.Interval == "daily" {
			return &domain.NotificationBroadcast{
				Platform: m.Platform,
				Title:    m.Title,
				LinkUrl:  m.LinkUrl,
				Time:     extractFristTime(m.FirstTime),
			}
		} else if *m.Interval == "weekly" {
			if repo.isBroadcastDay(m.FirstTime) {
				return &domain.NotificationBroadcast{
					Platform: m.Platform,
					Title:    m.Title,
					LinkUrl:  m.LinkUrl,
					Time:     extractFristTime(m.FirstTime),
				}
			}
		}
	case "pre":
		return &domain.NotificationBroadcast{
			Platform: m.Platform,
			Title:    m.Title,
			LinkUrl:  m.LinkUrl,
			Time:     extractFristTime(m.FirstTime),
		}
	default:
		return &domain.NotificationBroadcast{
			Platform: m.Platform,
			Title:    m.Title,
			LinkUrl:  m.LinkUrl,
			Time:     extractFristTime(m.FirstTime),
		}
	}

	return nil
}

func (repo *PlatformRepository) isBroadcastDay(ft *string) bool {
	firstDay, err := time.Parse("2006-01-02 15:04:05", *ft)
	if err != nil {
		return false
	}
	strToday := time.Now().Format("2006-01-02 15:04:05")
	today, err := time.Parse("2006-01-02 04:00:00", strToday)
	if err != nil {
		return false
	}

	diff := today.Sub(firstDay)
	diffHour := diff / time.Hour
	diffDay := int(diffHour / 24)

	if diffDay%7 == 0 {
		return true
	}
	return false
}

func extractFristTime(ft *string) *string {
	ftSlice := strings.Split(*ft, " ")
	return tools.NewNullString(ftSlice[1])
}
