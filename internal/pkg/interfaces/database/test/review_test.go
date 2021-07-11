package database_test

import (
	"animar/v1/internal/pkg/domain"
	infrastructure_test "animar/v1/internal/pkg/infrastructure/test"
	"animar/v1/internal/pkg/interfaces/database"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestReviewFilterAnime(t *testing.T) {
	/*   prepare   */
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("sqlmock not work")
	}
	defer db.Close()

	animeId := 1
	var content1, content2, content3 *string
	var rating1, rating2, rating3 *int
	rawContent1 := "よかった"
	rawContent3 := "神アニメ"
	rawRating2 := 8
	rawRating3 := 10
	userId1 := "aaaaaaaa"
	userId2 := "aaaaaaab"
	userId3 := "aaaaaaac"
	userId4 := ""
	content1 = &rawContent1
	content3 = &rawContent3
	rating2 = &rawRating2
	rating3 = &rawRating3

	rows := sqlmock.NewRows([]string{
		"id", "content", "rating", "anime_id", "user_id", "created_at", "updated_at",
	}).
		AddRow(1, content1, rating1, animeId, userId1, time.Now(), time.Now()).
		AddRow(2, content2, rating2, animeId, userId2, time.Now(), time.Now()).
		// AddRow(3, "良くなかった", 3, 2, "aaaaaaaa", time.Now(), time.Now()).
		AddRow(4, content3, rating3, animeId, userId3, time.Now(), time.Now())

	query := "Select * from reviews WHERE anime_id = ? AND (user_id != ? OR user_id IS NULL)"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(animeId, userId4).WillReturnRows(rows)

	repo := &database.ReviewRepository{
		SqlHandler: infrastructure_test.NewDummyHandler(db),
	}

	/*   prepare end   */
	/*   execute repository function   */

	reviews, err := repo.FilterByAnime(animeId, "")

	/*   execute repository function end   */
	/*   evaluation   */

	assert.Equal(t, err, nil)
	assert.Equal(t, len(reviews), 3)
	assert.Equal(t, reviews[0].Rating, rating1)
	assert.Equal(t, reviews[1].Content, content2)
	assert.Equal(t, reviews[2].Content, content3)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test GetAnimeReviews: %s", err)
	}
}

// update
func TestUpdateReviewContent(t *testing.T) {
	/*   prepare   */
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("sqlmock not work")
	}
	defer db.Close()

	rawContent := "面白かった。\n\n大満足!!"
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
	stmt.ExpectExec().WithArgs(rawContent, testReviewId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := &database.ReviewRepository{
		SqlHandler: infrastructure_test.NewDummyHandler(db),
	}
	r1 := &domain.TReviewInput{
		AnimeId: animeId,
		Content: rawContent,
	}

	/*   prepare end   */
	/*   execute repository function   */

	insertContent, err := repo.UpsertContent(*r1, userId)

	/*   execute repository function end   */
	/*   evaluation   */

	assert.Equal(t, err, nil)
	assert.Equal(t, insertContent, rawContent)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test InsertReviewContent: %s", err)
	}
}

// insert
func TestInsertReviewContent(t *testing.T) {
	/*   prepare   */
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("sqlmock not work")
	}
	defer db.Close()

	animeId := 1
	userId := "aaaaaaaa"
	content := "面白かった。\n\n大満足!!"

	rows := sqlmock.NewRows([]string{
		"id", "content", "rating", "anime_id", "user_id", "created_at", "updated_at",
	})

	query1 := "SELECT * FROM reviews WHERE anime_id = ? AND user_id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query1)).WithArgs(animeId, userId).WillReturnRows(rows)

	query2 := "INSERT INTO reviews(anime_id, content, user_id) VALUES(?, ?, ?)"
	stmt := mock.ExpectPrepare(regexp.QuoteMeta(query2))
	stmt.ExpectExec().WithArgs(animeId, content, userId).
		WillReturnResult(sqlmock.NewResult(20, 1))

	repo := &database.ReviewRepository{
		SqlHandler: infrastructure_test.NewDummyHandler(db),
	}
	r1 := &domain.TReviewInput{
		AnimeId: 1,
		Content: content,
	}

	/*   prepare end   */
	/*   execute repository function   */

	insertContent, err := repo.UpsertContent(*r1, userId)

	/*   execute repository function end   */
	/*   evaluation   */

	assert.Equal(t, err, nil)
	assert.Equal(t, insertContent, content)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Test InsertReviewContent: %s", err)
	}
}
