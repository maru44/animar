package database_test

import (
	"animar/v1/pkg/domain"
	infrastructure_test "animar/v1/pkg/infrastructure/test"
	"animar/v1/pkg/interfaces/database"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchAnimeReviews(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("sqlmock not work")
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id", "content", "rating", "anime_id", "user_id", "created_at", "updated_at",
	}).
		AddRow(1, "よかった", nil, 1, "aaaaaaaa", time.Now(), time.Now()).
		AddRow(2, nil, 8, 1, "aaaaaaab", time.Now(), time.Now()).
		// AddRow(3, "良くなかった", 3, 2, "aaaaaaaa", time.Now(), time.Now()).
		AddRow(4, "神アニメ", 10, 1, "aaaaaaac", time.Now(), time.Now())

	query := "Select * from reviews WHERE anime_id = ? AND (user_id != ? OR user_id IS NULL)"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(1, "").WillReturnRows(rows)

	repo := &database.ReviewRepository{
		SqlHandler: infrastructure_test.NewDummyHandler(db),
	}
	reviews, err := repo.FilterByAnime(1, "")

	var tempRating *int
	var tempContent *string

	assert.Equal(t, err, nil)
	assert.Equal(t, len(reviews), 3)
	tempRating = nil
	assert.Equal(t, reviews[0].Rating, tempRating)
	tempContent = nil
	assert.Equal(t, reviews[1].Content, tempContent)
	assert.Equal(t, *reviews[2].Content, "神アニメ")

	if err != nil {
		t.Errorf("FETCH FAIL: %s", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test GetAnimeReviews: %s", err)
	}
}

// update
func TestUpdateReviewContent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("sqlmock not work")
	}
	defer db.Close()

	animeId := 1
	userId := "aaaaaaaa"
	testReviewId := 1

	rows := sqlmock.NewRows([]string{
		"id", "content", "rating", "anime_id", "user_id", "created_at", "updated_at",
	}).
		AddRow(testReviewId, "よかった", 6, animeId, userId, time.Now(), time.Now())

	query1 := "SELECT * FROM reviews WHERE anime_id = ? AND user_id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query1)).WithArgs(animeId, userId).WillReturnRows(rows)

	query2 := "UPDATE reviews SET content = ? WHERE id = ?" // 対象がある場合のUPDATE
	stmt := mock.ExpectPrepare(regexp.QuoteMeta(query2))
	stmt.ExpectExec().WithArgs("面白かった。\n\n大満足!!", testReviewId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &database.ReviewRepository{
		SqlHandler: infrastructure_test.NewDummyHandler(db),
	}
	r1 := &domain.TReviewInput{
		AnimeId: 1,
		Content: "面白かった。\n\n大満足!!",
	}
	insertContent, err := repo.UpsertContent(*r1, userId)

	assert.Equal(t, err, nil)
	assert.Equal(t, insertContent, "面白かった。\n\n大満足!!")

	if err != nil {
		t.Errorf("INSERT Review content FAIL: %s", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test InsertReviewContent: %s", err)
	}
}

// insert
func TestInsertReviewContent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("sqlmock not work")
	}
	defer db.Close()

	animeId := 1
	userId := "aaaaaaaa"

	rows := sqlmock.NewRows([]string{
		"id", "content", "rating", "anime_id", "user_id", "created_at", "updated_at",
	}).
		AddRow(0, "", 0, 0, "", time.Now(), time.Now())

	query1 := "SELECT * FROM reviews WHERE anime_id = ? AND user_id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query1)).WithArgs(animeId, userId).WillReturnRows(rows)

	query2 := "INSERT INTO reviews(anime_id, content, user_id) VALUES(?, ?, ?)"
	stmt := mock.ExpectPrepare(regexp.QuoteMeta(query2))
	stmt.ExpectExec().WithArgs(animeId, "面白かった。\n\n大満足!!", userId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &database.ReviewRepository{
		SqlHandler: infrastructure_test.NewDummyHandler(db),
	}
	r1 := &domain.TReviewInput{
		AnimeId: 1,
		Content: "面白かった。\n\n大満足!!",
	}
	insertContent, err := repo.UpsertContent(*r1, userId)

	assert.Equal(t, err, nil)
	assert.Equal(t, insertContent, "面白かった。\n\n大満足!!")

	if err != nil {
		t.Errorf("INSERT Review content FAIL: %s", err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test InsertReviewContent: %s", err)
	}
}
