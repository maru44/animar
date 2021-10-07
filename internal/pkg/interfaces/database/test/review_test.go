package database_test

import (
	"animar/v1/internal/pkg/domain"
	infrastructure_test "animar/v1/internal/pkg/infrastructure/test"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/tools/tools"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/maru44/perr"
	"github.com/stretchr/testify/assert"
)

func Test_FilterByAnime(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &database.ReviewRepository{
		SqlHandler: infrastructure_test.NewDummyHandler(db),
	}

	table := []struct {
		testName string
		animeId  int
		userId   string
		reviews  domain.TReviews
		err      error
	}{
		{
			"success 1",
			24,
			"user_no_054",
			domain.TReviews{
				{
					ID:        800,
					Content:   tools.NewNullString("素晴らしいアニメダス"),
					Rating:    nil,
					AnimeId:   24,
					UserId:    tools.NewNullString("user_no_754"),
					CreatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
				},
				{
					ID:        933,
					Content:   nil,
					Rating:    tools.NewNullInt(9),
					AnimeId:   24,
					UserId:    tools.NewNullString("user_akwaejjf"),
					CreatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
					UpdatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
				},
			},
			nil,
		},
	}

	q := "Select * from reviews WHERE anime_id = ? AND (user_id != ? OR user_id IS NULL)"

	for _, tt := range table {
		t.Run(tt.testName, func(t *testing.T) {
			rows := sqlmock.NewRows([]string{
				"reviews_id", "reviews_content", "reviews_rating",
				"reviews_anime_id", "reviews_user_id",
				"reviews_created_at", "reviews_updated_at",
			})
			for _, r := range tt.reviews {
				rows.AddRow(r.ID, r.Content, r.Rating, r.AnimeId, r.UserId, r.CreatedAt, r.UpdatedAt)
			}

			mock.ExpectQuery(regexp.QuoteMeta(q)).WithArgs(tt.animeId, tt.userId).WillReturnRows(rows)

			rs, err := repo.FilterByAnime(tt.animeId, tt.userId)

			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.reviews, rs)
		})
	}
}

func Test_UpsertContent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := &database.ReviewRepository{
		SqlHandler: infrastructure_test.NewDummyHandler(db),
	}

	q1 := "UPDATE reviews SET content = ? WHERE anime_id = ? AND user_id = ?"
	q2 := "INSERT INTO reviews(anime_id, content, user_id) VALUES(?, ?, ?)"

	table := []struct {
		name         string
		ri           domain.TReviewInput
		userId       string
		affected     int64
		lastInserted int64
		retContent   string
		noError      bool
	}{
		{
			"success update",
			domain.TReviewInput{
				AnimeId: 122,
				Content: "つまらないと思ったけど実は面白かった",
				Rating:  0,
			},
			"user+desuyo",
			1,
			0,
			"つまらないと思ったけど実は面白かった",
			true,
		},
		{
			"success insert",
			domain.TReviewInput{
				AnimeId: 134,
				Content: "初見 最高でした",
				Rating:  0,
			},
			"user+desuyo",
			0,
			64,
			"初見 最高でした",
			true,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			stmt1 := mock.ExpectPrepare(regexp.QuoteMeta(q1))
			stmt1.ExpectExec().WithArgs(tt.ri.Content, tt.ri.AnimeId, tt.userId).
				WillReturnResult(sqlmock.NewResult(0, tt.affected))

			if tt.affected == 0 {
				stmt2 := mock.ExpectPrepare(regexp.QuoteMeta(q2))
				stmt2.ExpectExec().WithArgs(tt.ri.AnimeId, tt.ri.Content, tt.userId).
					WillReturnResult(sqlmock.NewResult(tt.lastInserted, 0))
			}

			result, err := repo.UpsertContent(tt.ri, tt.userId)

			assert.Equal(t, tt.noError, perr.IsNoError(err))
			assert.Equal(t, tt.retContent, result)
		})
	}
}
