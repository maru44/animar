package review

import (
	"animar/v1/auth"
	"context"
)

func OnesReviewsDomain(userId string) []TReview {
	rows := OnesReviewsList(userId)
	var revs []TReview
	for rows.Next() {
		var rev TReview
		err := rows.Scan(&rev.ID, &rev.Content, &rev.Star, &rev.AnimeId, &rev.UserId, &rev.CreatedAt, &rev.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		revs = append(revs, rev)
	}
	defer rows.Close()
	return revs
}

// user以外のレビュー
func AnimeReviewsDomain(animeId int, userId string) []TReview {
	rows := AnimeReviewsList(animeId, userId)
	var revs []TReview
	for rows.Next() {
		var rev TReview
		err := rows.Scan(&rev.ID, &rev.Content, &rev.Star, &rev.AnimeId, &rev.UserId, &rev.CreatedAt, &rev.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		revs = append(revs, rev)
	}

	defer rows.Close()
	return revs
}

func AnimeReviewsWithUserInfoDomain(animeId int, userId string) []TReviewJoinUser {
	ctx := context.Background()
	rows := AnimeReviewsList(animeId, userId)
	var reviews []TReviewJoinUser
	for rows.Next() {
		var rev TReviewJoinUser
		err := rows.Scan(
			&rev.ID, &rev.Content, &rev.Star, &rev.AnimeId, &rev.UserId, &rev.CreatedAt, &rev.UpdatedAt,
		)
		if rev.UserId != nil {
			user := auth.GetUserFirebase(ctx, *rev.UserId)
			rev.User = user
		} else {
			rev.User = nil
		}

		if err != nil {
			continue
			//fmt.Print(err)
		}
		reviews = append(reviews, rev)
	}
	defer rows.Close()
	return reviews
}

func OnesReviewsJoinAnimeDomain(userId string) []TReviewJoinAnime {
	rows := OnesReviewsJoinAnime(userId)
	var revs []TReviewJoinAnime
	for rows.Next() {
		var rev TReviewJoinAnime
		err := rows.Scan(
			&rev.ID, &rev.Content, &rev.Star, &rev.AnimeId, &rev.UserId,
			&rev.CreatedAt, &rev.UpdatedAt, &rev.Title, &rev.Slug, &rev.AnimeContent, &rev.OnAirState,
		)
		if err != nil {
			panic(err.Error())
		}
		revs = append(revs, rev)
	}
	defer rows.Close()
	return revs
}
