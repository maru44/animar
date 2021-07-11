package usecase_test

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/infrastructure"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/tools/tools"
	"animar/v1/internal/pkg/usecase"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newSqlHandler(sql *sql.DB) database.SqlHandler {
	sqlHandler := new(infrastructure.SqlHandler)
	sqlHandler.Conn = sql
	return sqlHandler
}

// insert review
func TestInsertReviewContent(t *testing.T) {
	r := &domain.TReviewInput{
		AnimeId: 1,
		Content: "感想 一言\nいいね",
		Rating:  8,
	}
	userId := "aaaaaaaa"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("DB FAIL: %s", err)
	}
	query := "INSERT INTO reviews(anime_id, content, user_id) VALUES(?, ?, ?)"
	stmt := mock.ExpectPrepare(regexp.QuoteMeta(query))
	stmt.ExpectExec().WithArgs(r.AnimeId, r.Content, userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlHandler := newSqlHandler(db)
	us := usecase.NewReviewInteractor(
		&database.ReviewRepository{
			SqlHandler: sqlHandler,
		},
	)
	_, err = us.PostReviewContent(*r, userId)
	if err != nil {
		t.Errorf("INSERT FAIL: %s", err)
	}
}

// animes review
func TestFetchAnimeReviews(t *testing.T) {
	now := time.Now()
	nowStr := now.String()
	m := domain.TReviews{
		domain.TReview{
			ID: 1, Content: tools.NewNullString("最高のアニメ!!"), Rating: tools.NewNullInt(8), AnimeId: 12,
			UserId: tools.NewNullString("a"), CreatedAt: nowStr, UpdatedAt: nowStr,
		},
		domain.TReview{
			ID: 2, Content: tools.NewNullString("神アニメ!!"), Rating: tools.NewNullInt(10), AnimeId: 12,
			UserId: tools.NewNullString("b"), CreatedAt: nowStr, UpdatedAt: nowStr,
		},
		domain.TReview{
			ID: 3, Content: tools.NewNullString("普通"), Rating: tools.NewNullInt(6),
			AnimeId: 13, UserId: tools.NewNullString("b"), CreatedAt: nowStr, UpdatedAt: nowStr,
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("DB FAIL: %s", err)
	}
	rows := sqlmock.NewRows([]string{"id", "rating", "anime_id", "user_id", "created_at", "updated_at"}).
		AddRow(m[0].ID, m[0].Content, m[0].Rating, m[0].AnimeId, m[0].UserId, m[0].CreatedAt, m[0].UpdatedAt).
		AddRow(m[1].ID, m[1].Content, m[1].Rating, m[1].AnimeId, m[1].UserId, m[1].CreatedAt, m[1].UpdatedAt)

	query := "Select * from reviews WHERE anime_id = ? AND (user_id != ? OR user_id IS NULL)"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(12, "").
		WillReturnRows(rows)

	sqlHandler := newSqlHandler(db)
	r := usecase.NewReviewInteractor(
		&database.ReviewRepository{
			SqlHandler: sqlHandler,
		},
	)
	revAnime, err := r.GetAnimeReviews(12, "")

	assert.Equal(t, err, nil)
	assert.Equal(t, len(revAnime), 2)

	if err != nil {
		t.Errorf("INSERT FAIL: %s", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test GetAnimeReviews: %s", err)
	}
}
