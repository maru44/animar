package review

import (
	"animar/v1/helper"
	"database/sql"
)

type TReview struct {
	ID        int
	Content   string
	Star      int
	AnimeId   int
	UserId    string
	CreatedAt string
	UpdatedAt string
}

func ListReviews() *sql.Rows {
	db := helper.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from tbl_review")
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func DetailReview(id int) TReview {
	db := helper.AccessDB()
	defer db.Close()

	var rev TReview
	nullStar := new(helper.NullInt)
	nullContent := new(sql.NullString)
	nullUserId := new(sql.NullString)
	err := db.QueryRow("SELECT * FROM tbl_reviews WHERE id = ?", id).Scan(&rev.ID, nullStar, nullContent, &rev.AnimeId, nullUserId, &rev.CreatedAt, &rev.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		rev.ID = 0
	case err != nil:
		panic(err.Error())
	default:
		rev.Content = nullContent.String
		rev.Star = *nullStar.Int
		rev.UserId = nullUserId.String
	}
	return rev
}

func OnesReviewsList(userId string) *sql.Rows {
	db := helper.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from tbl_reviews WHERE user_id = ?", userId)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func AnimeReviewsList(animeId int) *sql.Rows {
	db := helper.AccessDB()
	defer db.Close()
	rows, err := db.Query("Select * from tbl_reviews WHERE anime_id = ?", animeId)
	if err != nil {
		panic(err.Error())
	}
	return rows
}

func InsertReview(animeId int, content string, star int, user_id string) int {
	db := helper.AccessDB()
	defer db.Close()

	stmtInsert, err := db.Prepare("INSERT INTO tbl_reviews(anime_id, content, star, user_id) VALUES(?, ?, ?, ?)")
	defer stmtInsert.Close()

	exe, err := stmtInsert.Exec(animeId, helper.NewNullString(content).String, helper.NewNullInt(star).Int, user_id)

	insertedId, err := exe.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	return int(insertedId)
}
