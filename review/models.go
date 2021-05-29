package review

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

// nullableはpointerにしてnilを受け取れるようにする
type TReview struct {
	ID        int
	Content   *string
	Star      *int
	AnimeId   int
	UserId    *string
	CreatedAt string
	UpdatedAt string
}

type TReviewJoinAnime struct {
	ID           int
	Content      *string
	Star         *int
	AnimeId      int
	UserId       *string
	CreatedAt    string
	UpdatedAt    string
	Title        string
	Slug         string
	AnimeContent *string
	OnAirState   *int
}

type TReviewJoinUser struct {
	ID        int
	Content   *string
	Star      *int
	AnimeId   int
	UserId    *string
	CreatedAt string
	UpdatedAt string
	User      *auth.UserInfo
}

// all
func ListReviews() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from tbl_reviews")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

// retrieve
// by id
func DetailReview(id int) TReview {
	db := tools.AccessDB()
	defer db.Close()

	var rev TReview
	err := db.QueryRow("SELECT * FROM tbl_reviews WHERE id = ?", id).Scan(&rev.ID, &rev.Star, &rev.Content, &rev.AnimeId, &rev.UserId, &rev.CreatedAt, &rev.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		rev.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return rev
}

// ones detail review
func DetailReviewAnimeUser(animeId int, userId string) TReview {
	db := tools.AccessDB()
	defer db.Close()

	var rev TReview
	err := db.QueryRow(
		"SELECT * FROM tbl_reviews WHERE anime_id = ? AND user_id = ?", animeId, userId,
	).Scan(
		&rev.ID, &rev.Content, &rev.Star, &rev.AnimeId, &rev.UserId, &rev.CreatedAt, &rev.UpdatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		rev.ID = 0
	case err != nil:
		panic(err.Error())
	}
	return rev
}

// List
// fiter by userId
func OnesReviewsList(userId string) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from tbl_reviews WHERE user_id = ?", userId)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

// filter by userId
// joined anime detail
func OnesReviewsJoinAnime(userId string) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query(
		"SELECT tbl_reviews.*, anime.title, anime.slug, anime.content, anime.on_air_state "+
			"FROM tbl_reviews LEFT JOIN anime ON tbl_reviews.anime_id = anime.id WHERE user_id = ?", userId,
	)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

// @TODO exclude star
// List
// filter by animeId
func AnimeReviewsList(animeId int, userId string) *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from tbl_reviews WHERE anime_id = ? AND (user_id != ? OR user_id IS NULL)", animeId, userId)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

// Post
func InsertReview(animeId int, content string, star int, user_id string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO tbl_reviews(anime_id, content, star, user_id) VALUES(?, ?, ?, ?)")
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(animeId, tools.NewNullString(content).String, tools.NewNullInt(star).Int, user_id)

	insertedId, err := exe.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	return int(insertedId)
}

// insert content of review
func InsertReviewContent(animeId int, userId string, content string) string {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO tbl_reviews(anime_id, content, user_id) VALUES(?, ?, ?)")
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(animeId, tools.NewNullString(content).String, userId)
	insertedId, err := exe.LastInsertId()
	fmt.Print(int(insertedId))
	if err != nil {
		panic(err.Error())
	}

	return content
}

// upsert content of review
func UpsertReviewContent(animeId int, content string, userId string) string {
	db := tools.AccessDB()
	defer db.Close()

	var rev TReview
	err := db.QueryRow("SELECT id from tbl_reviews WHERE user_id = ? AND anime_id = ?",
		userId, animeId).Scan(&rev.ID)

	switch {
	case err == sql.ErrNoRows:
		return InsertReviewContent(animeId, userId, content)
	case err != nil:
		panic(err.Error())
	default:
		// update
		stmtUpdate, err := db.Prepare("UPDATE tbl_reviews SET content = ? WHERE id = ?")
		defer stmtUpdate.Close()
		exe, err := stmtUpdate.Exec(content, &rev.ID)
		fmt.Print(exe)
		if err != nil {
			panic(err.Error())
		}
		return content
	}
}

// insert star of review
func InsertReviewStar(animeId int, userId string, star int) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO tbl_reviews(anime_id, star, user_id) VALUES(?, ?, ?)")
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(animeId, tools.NewNullInt(star).Int, userId)
	insertedId, err := exe.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	fmt.Print(insertedId)
	return star
}

// upsert star of content
func UpsertReviewStar(animeId int, star int, userId string) int {
	db := tools.AccessDB()
	defer db.Close()

	var rev TReview
	err := db.QueryRow("SELECT id from tbl_reviews WHERE user_id = ? AND anime_id = ?",
		userId, animeId).Scan(&rev.ID)

	switch {
	case err == sql.ErrNoRows:
		return InsertReviewStar(animeId, userId, star)
	case err != nil:
		panic(err.Error())
	default:
		// update
		stmtUpdate, err := db.Prepare("UPDATE tbl_reviews SET star = ? WHERE id = ?")
		defer stmtUpdate.Close()
		exe, err := stmtUpdate.Exec(star, &rev.ID)
		if err != nil {
			panic(err.Error())
		}
		//insertedId, _ := exe.RowsAffected()
		fmt.Print(exe)
		return star
	}
}

func AnimeStarAvg(animeId int) string {
	db := tools.AccessDB()
	defer db.Close()

	var avg float32
	err := db.QueryRow("SELECT COALESCE(AVG(star), 0) FROM tbl_reviews WHERE anime_id = ?", animeId).Scan(&avg)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%.1f", avg)
}
