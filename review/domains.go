package review

import (
	"animar/v1/helper"
	"database/sql"
)

func OnesReviewsDomain(userId string) []TReview {
	rows := OnesReviewsList(userId)
	var revs []TReview
	for rows.Next() {
		var rev TReview
		nullContent := new(sql.NullString)
		nullStar := new(helper.NullInt)
		nullUserId := new(sql.NullString)
		err := rows.Scan(&rev.ID, nullContent, nullStar, &rev.AnimeId, nullUserId, &rev.CreatedAt, &rev.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		rev.Content = nullContent.String
		rev.Star = nullStar.Int
		rev.UserId = nullContent.String
		revs = append(revs, rev)
	}

	defer rows.Close()

	return revs
}
