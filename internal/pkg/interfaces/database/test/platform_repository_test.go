package database_test

import (
	"animar/v1/internal/pkg/domain"
	infrastructure_test "animar/v1/internal/pkg/infrastructure/test"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/interfaces/database/queryset"
	"animar/v1/internal/pkg/tools/tools"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPlatform_FilterTodaysBroadcast(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &database.PlatformRepository{
		SqlHandler: infrastructure_test.NewDummyHandler(db),
	}

	table := []struct {
		testName      string
		rawBroadCasts []domain.RawNotificationMaterial // result of Query
		broadcasts    []domain.NotificationBroadcast   // result
		err           error
	}{
		{
			"now filtering",
			[]domain.RawNotificationMaterial{
				{
					Platform:  "netflix",
					Title:     "宇宙よりも遠い場所",
					LinkUrl:   tools.NewNullString("https://abc.example.com"),
					BaseUrl:   tools.NewNullString("https://netflix.com"),
					FirstTime: tools.NewNullString(time.Now().Add(-time.Hour*24*7*4).Format("2006-01-02") + " 04:01:00"),
					Interval:  tools.NewNullString("weekly"),
					State:     "now",
				},
				{
					Platform:  "amapri",
					Title:     "ガンダム",
					LinkUrl:   nil,
					BaseUrl:   tools.NewNullString("https://amazon.com"),
					FirstTime: tools.NewNullString(time.Now().Add(time.Hour*24).Format("2006-01-02") + " 04:00:00"),
					Interval:  tools.NewNullString("weekly"),
					State:     "now",
				},
				{
					Platform:  "dani",
					Title:     "攻殻機動隊S.A.C.",
					LinkUrl:   nil,
					BaseUrl:   tools.NewNullString("https://danime.com"),
					FirstTime: tools.NewNullString(time.Now().Add(-time.Hour*24*7).Format("2006-01-02") + " 03:59:00"),
					Interval:  tools.NewNullString("weekly"),
					State:     "now",
				},
				{
					Platform:  "dani",
					Title:     "おジャ魔女ドレミ",
					LinkUrl:   tools.NewNullString("https://oja.majo"),
					BaseUrl:   tools.NewNullString("https://danime.com"),
					FirstTime: tools.NewNullString(time.Now().Add(time.Hour*24).Format("2006-01-02") + " 14:01:00"),
					Interval:  tools.NewNullString("weekly"),
					State:     "now",
				},
				{
					Platform:  "dani",
					Title:     "おジャ魔女ドレミ#",
					LinkUrl:   tools.NewNullString("https://oja.majo"),
					BaseUrl:   tools.NewNullString("https://danime.com"),
					FirstTime: tools.NewNullString(time.Now().Add(time.Hour*24*2).Format("2006-01-02") + " 05:01:00"),
					Interval:  tools.NewNullString("weekly"),
					State:     "now",
				},
			},
			[]domain.NotificationBroadcast{
				{
					Platform: "netflix",
					Title:    "宇宙よりも遠い場所",
					LinkUrl:  tools.NewNullString("https://abc.example.com"),
					Time:     tools.NewNullString("04:01:00"),
				},
				{
					Platform: "amapri",
					Title:    "ガンダム",
					LinkUrl:  tools.NewNullString("https://amazon.com"),
					Time:     tools.NewNullString("04:00:00"),
				},
			},
			nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.testName, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{
				"plat_plat_name", "animes_title", "rel_link_url",
				"plat_base_url", "rel_first_broadcast", "rel_delivery_interval", "animes_state",
			})
			for _, rb := range tt.rawBroadCasts {
				fmt.Println(*rb.FirstTime)
				rows.AddRow(
					rb.Platform, rb.Title, rb.LinkUrl, rb.BaseUrl,
					rb.FirstTime, rb.Interval, rb.State,
				)
			}
			mock.ExpectQuery(regexp.QuoteMeta(queryset.TodaysBroadcastQuery)).WillReturnRows(rows)

			bs, err := repo.FilterTodaysBroadCast()

			if err != nil {
				fmt.Println(err.Error())
			}

			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.broadcasts, bs)
		})
	}
}
