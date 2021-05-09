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
	err := db.QueryRow("SELECT * FROM tbl_reviews WHERE id = ?", id).Scan(&rev.ID, nullStar, nullContent, &rev.AnimeId, &rev.CreatedAt, &rev.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		rev.ID = 0
	case err != nil:
		panic(err.Error())
	default:
		rev.Content = nullContent.String
		rev.Star = nullStar.Int
	}
	return rev
}