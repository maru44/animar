package database

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database/queryset"
	"animar/v1/internal/pkg/tools/tools"
	"fmt"
	"strings"
	"time"

	"github.com/maru44/perr"
)

type PlatformRepository struct {
	SqlHandler
}

func (repo *PlatformRepository) FilterByAnime(animeId int) (platforms domain.TRelationPlatforms, err error) {
	rows, err := repo.Query(queryset.PlatformFilterByAnimeQuery, animeId)
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
	rows, err := repo.Query(queryset.TodaysBroadcastQuery)
	if err != nil {
		return nil, perr.Wrap(err, perr.BadRequest)
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

		b, err := repo.modifyBroadcastTime(&m)
		if err != nil {
			return nil, perr.Wrap(err, perr.BadRequest)
		}
		if b != nil {
			out = append(out, *b)
		}
	}

	return out, nil
}

func (repo *PlatformRepository) modifyBroadcastTime(m *domain.RawNotificationMaterial) (*domain.NotificationBroadcast, error) {
	if m.Interval == nil {
		return nil, nil
	}

	link := m.LinkUrl
	if link == nil {
		link = m.BaseUrl
	}

	switch m.State {
	case "now":
		if *m.Interval == "once" || *m.Interval == "daily" {
			return &domain.NotificationBroadcast{
				Platform: m.Platform,
				Title:    m.Title,
				LinkUrl:  link,
				Time:     extractFristTime(m.FirstTime),
			}, nil
		} else if *m.Interval == "weekly" {
			if b, err := repo.isBroadcastDay(m.FirstTime); err != nil {
				return nil, perr.Wrap(err, perr.BadRequest)
			} else {
				if b {
					return &domain.NotificationBroadcast{
						Platform: m.Platform,
						Title:    m.Title,
						LinkUrl:  link,
						Time:     extractFristTime(m.FirstTime),
					}, nil
				}
			}
		}
	case "pre":
		return &domain.NotificationBroadcast{
			Platform: m.Platform,
			Title:    m.Title,
			LinkUrl:  link,
			Time:     extractFristTime(m.FirstTime),
		}, nil
	default:
		return &domain.NotificationBroadcast{
			Platform: m.Platform,
			Title:    m.Title,
			LinkUrl:  link,
			Time:     extractFristTime(m.FirstTime),
		}, nil
	}

	return nil, nil
}

func (repo *PlatformRepository) isBroadcastDay(ft *string) (bool, error) {
	firstDay, err := time.Parse("2006-01-02 15:04:05 MST", *ft+" JST")
	if err != nil {
		return false, perr.Wrap(err, perr.BadRequest)
	}

	strToday := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	today, err := time.Parse("2006-01-02 15:04:05 MST", strToday+" 04:00:00 JST")
	if err != nil {
		return false, perr.Wrap(err, perr.BadRequest)
	}

	diff := today.Sub(firstDay)
	diffHour := diff / time.Hour
	diffDay := int(diffHour / 24)
	fmt.Println(diffDay)

	if diffDay%7 == 0 {
		return true, nil
	}
	return false, nil
}

func extractFristTime(ft *string) *string {
	ftSlice := strings.Split(*ft, " ")
	return tools.NewNullString(ftSlice[1])
}
