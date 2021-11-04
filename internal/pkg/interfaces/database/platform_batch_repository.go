package database

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database/queryset"
	"animar/v1/internal/pkg/tools/tools"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/maru44/perr"
)

type PlatformBatchRepository struct {
	SqlHandler
}

func (repo *PlatformBatchRepository) FilterNotificationTarget() ([]string, error) {
	rows, err := repo.Query(queryset.NotificationTarget)
	if err != nil {
		return nil, perr.Wrap(err, perr.BadRequest)
	}

	rows.Next()
	var targetsStr string
	err = rows.Scan(&targetsStr)
	if err != nil {
		return nil, perr.Wrap(err, perr.BadRequest)
	}

	targets := strings.Split(targetsStr, ",")
	return targets, nil
}

func (repo *PlatformBatchRepository) FilterTodaysBroadCast() ([]domain.NotificationBroadcast, error) {
	rows, err := repo.Query(queryset.TodaysBroadcastQuery)
	if err != nil {
		return nil, perr.Wrap(err, perr.BadRequest)
	}

	var out []domain.NotificationBroadcast
	for rows.Next() {
		var m domain.RawNotificationMaterial
		err := rows.Scan(
			&m.Platform, &m.Title, &m.Slug, &m.LinkUrl, &m.BaseUrl,
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

func (repo *PlatformBatchRepository) MakeSlackMessage(nbs []domain.NotificationBroadcast) (out string) {
	animesMap := map[string]string{}
	for _, n := range nbs {
		if n.Time != nil && n.LinkUrl != nil {
			if _, ok := animesMap[n.Slug]; ok {
				animesMap[n.Slug] += fmt.Sprintf(
					"\t◦ %s %s <%s|:link:> <https://%s%s/anime/%s|:heart:>\n",
					n.Platform, *n.Time, *n.LinkUrl, configs.FrontHost, configs.FrontPort, n.Slug,
				)
			} else {
				animesMap[n.Slug] = fmt.Sprintf(
					"• %s\n\t◦ %s %s <%s|:link:> <https://%s%s/anime/%s|:heart:>\n",
					n.Title, n.Platform, *n.Time, *n.LinkUrl, configs.FrontHost, configs.FrontPort, n.Slug,
				)
			}
		}
	}

	for _, s := range animesMap {
		out += s
	}
	return out
}

func (repo *PlatformBatchRepository) ChangeOnAirState(state string, slug string) (int, error) {
	exe, err := repo.Execute(
		"UPDATE animes SET state = ? WHERE slug = ?",
		state, slug,
	)
	if err != nil {
		return 0, perr.Wrap(err, perr.BadRequest)
	}

	rawAffected, err := exe.LastInsertId()
	if err != nil {
		return 0, perr.Wrap(err, perr.BadRequest)
	}
	return int(rawAffected), nil
}

func (repo *PlatformBatchRepository) modifyBroadcastTime(m *domain.RawNotificationMaterial) (*domain.NotificationBroadcast, error) {
	if m.Interval == nil {
		return nil, nil
	}

	link := m.LinkUrl
	if link == nil {
		link = m.BaseUrl
	}

	out := &domain.NotificationBroadcast{
		Platform: m.Platform,
		Title:    m.Title,
		Slug:     m.Slug,
		LinkUrl:  link,
		Time:     extractFristTime(m.FirstTime),
	}

	switch m.State {
	case domain.StateNow:
		if *m.Interval == "once" || *m.Interval == "daily" {
			return out, nil
		} else if *m.Interval == "weekly" {
			if b, err := repo.isBroadcastDay(m.FirstTime); err != nil {
				return nil, perr.Wrap(err, perr.BadRequest)
			} else {
				if b {
					return out, nil
				}
			}
		}
	case domain.StatePre:
		repo.ChangeOnAirState(domain.StateNow, out.Slug)
		return out, nil
	default:
		return out, nil
	}

	return nil, nil
}

func (repo *PlatformBatchRepository) isBroadcastDay(ft *string) (bool, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return false, perr.Wrap(err, perr.BadRequest)
	}
	firstDay, err := time.Parse("2006-01-02 15:04:05", *ft)
	if err != nil {
		return false, perr.Wrap(err, perr.BadRequest)
	}
	firstDay = firstDay.In(jst)

	strToday := time.Now().Add(time.Hour * 24).Format("2006-01-02")
	today, err := time.Parse("2006-01-02 15:04:05", strToday+" 04:00:00")
	if err != nil {
		return false, perr.Wrap(err, perr.BadRequest)
	}
	today = today.In(jst)

	diff := today.Sub(firstDay)
	if diff < 0 {
		return false, nil
	}
	diffHour := int(math.Floor(float64(diff / time.Hour)))
	diffDay := diffHour / 24

	if diffDay%7 == 0 {
		return true, nil
	}
	return false, nil
}

func extractFristTime(ft *string) *string {
	if ft == nil {
		return nil
	}
	ftSlice := strings.Split(*ft, " ")
	t := ftSlice[1]
	// 00 ~ 03 hour  >> 24 ~ 27
	hour := string([]rune(t)[:2])
	min := string([]rune(t)[2:5])
	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		log.Print(err)
	}
	if hourInt < 4 {
		hourInt += 24
	}
	return tools.NewNullString(strconv.Itoa(hourInt) + min)
}
