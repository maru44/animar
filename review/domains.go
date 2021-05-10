package review

import (
	"database/sql"
)

func OnesReviewsDomain(userId string) []TReview {
	rows := OnesReviewsList(userId)
	var revs []TReview
	for rows.Next() {
		var rev TReview
		nullContent := new(sql.NullString)
		nullStar := new(sql.NullInt32)
		nullUserId := new(sql.NullString)
		err := rows.Scan(&rev.ID, nullContent, nullStar, &rev.AnimeId, nullUserId, &rev.CreatedAt, &rev.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		rev.Content = nullContent.String
		rev.Star = int(nullStar.Int32)
		rev.UserId = nullUserId.String
		revs = append(revs, rev)
	}

	defer rows.Close()

	return revs
}
