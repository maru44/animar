package review

import (
	"animar/v1/tools"
	"database/sql"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

// nullableはpointerにしてnilを受け取れるようにする
type TReview struct {
	ID        int     `json:"id"`
	Content   *string `json:"content,omitempty"`
	Rating    *int    `json:"rating,omitempty"`
	AnimeId   int     `json:"anime_id"`
	UserId    *string `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

type TReviewJoinAnime struct {
	ID           int     `json:"id"`
	Content      *string `json:"content,omitempty"`
	Rating       *int    `json:"rating,omitempty"`
	AnimeId      int     `json:"anime_id"`
	UserId       *string `json:"user_id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
	Title        string  `json:"title"`
	Slug         string  `json:"slug"`
	AnimeContent *string `json:"anime_content,omitempty"`
	AState       *int    `json:"anime_state,omitempty"`
}

type TReviewJoinUser struct {
	ID        int     `json:"id"`
	Content   *string `json:"content,omitempty"`
	Rating    *int    `json:"rating,omitempty"`
	AnimeId   int     `json:"anime_id"`
	UserId    *string `json:"user_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at,omitempty"`
	User      *auth.UserInfo
}

// all
func ListReviews() *sql.Rows {
	db := tools.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from reviews")
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
	err := db.QueryRow("SELECT * FROM reviews WHERE id = ?", id).Scan(&rev.ID, &rev.Rating, &rev.Content, &rev.AnimeId, &rev.UserId, &rev.CreatedAt, &rev.UpdatedAt)

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
		"SELECT * FROM reviews WHERE anime_id = ? AND user_id = ?", animeId, userId,
	).Scan(
		&rev.ID, &rev.Content, &rev.Rating, &rev.AnimeId, &rev.UserId, &rev.CreatedAt, &rev.UpdatedAt,
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
	rows, err := db.Query("Select * from reviews WHERE user_id = ?", userId)
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
		"SELECT reviews.*, animes.title, animes.slug, animes.content, animes.on_air_state "+
			"FROM reviews LEFT JOIN animes ON reviews.anime_id = animes.id WHERE user_id = ?", userId,
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
	rows, err := db.Query("Select * from reviews WHERE anime_id = ? AND (user_id != ? OR user_id IS NULL)", animeId, userId)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

// Post
func InsertReview(animeId int, content string, rating int, user_id string) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO reviews(anime_id, content, rating, user_id) VALUES(?, ?, ?, ?)")
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(animeId, content, rating, user_id)
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

	stmtInsert, err := db.Prepare("INSERT INTO reviews(anime_id, content, user_id) VALUES(?, ?, ?)")
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(animeId, content, userId)
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
	err := db.QueryRow("SELECT id from reviews WHERE user_id = ? AND anime_id = ?",
		userId, animeId).Scan(&rev.ID)

	switch {
	case err == sql.ErrNoRows:
		return InsertReviewContent(animeId, userId, content)
	case err != nil:
		panic(err.Error())
	default:
		// update
		stmtUpdate, err := db.Prepare("UPDATE reviews SET content = ? WHERE id = ?")
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
func InsertReviewStar(animeId int, userId string, rating int) int {
	db := tools.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO reviews(anime_id, rating, user_id) VALUES(?, ?, ?)")
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(animeId, rating, userId)
	insertedId, err := exe.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	fmt.Print(insertedId)
	return rating
}

// upsert star of content
func UpsertReviewStar(animeId int, rating int, userId string) int {
	db := tools.AccessDB()
	defer db.Close()

	var rev TReview
	err := db.QueryRow("SELECT id from reviews WHERE user_id = ? AND anime_id = ?",
		userId, animeId).Scan(&rev.ID)

	switch {
	case err == sql.ErrNoRows:
		return InsertReviewStar(animeId, userId, rating)
	case err != nil:
		panic(err.Error())
	default:
		// update
		stmtUpdate, err := db.Prepare("UPDATE reviews SET rating = ? WHERE id = ?")
		defer stmtUpdate.Close()
		exe, err := stmtUpdate.Exec(rating, &rev.ID)
		if err != nil {
			panic(err.Error())
		}
		//insertedId, _ := exe.RowsAffected()
		fmt.Print(exe)
		return rating
	}
}

func AnimeStarAvg(animeId int) string {
	db := tools.AccessDB()
	defer db.Close()

	var avg float32
	err := db.QueryRow("SELECT COALESCE(AVG(rating), 0) FROM reviews WHERE anime_id = ?", animeId).Scan(&avg)
	if err != nil {
		panic(err.Error())
	}
	return fmt.Sprintf("%.1f", avg)
}
