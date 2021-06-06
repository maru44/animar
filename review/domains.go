package review

import (
	"animar/v1/auth"
	"context"
)

func OnesReviewsDomain(userId string) []TReview {
	rows := OnesReviewsList(userId)
	var revs []TReview
	for rows.Next() {
		var r TReview
		err := rows.Scan(&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		revs = append(revs, r)
	}
	defer rows.Close()
	return revs
}

// user以外のレビュー
func AnimeReviewsDomain(animeId int, userId string) []TReview {
	rows := AnimeReviewsList(animeId, userId)
	var revs []TReview
	for rows.Next() {
		var r TReview
		err := rows.Scan(&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		revs = append(revs, r)
	}

	defer rows.Close()
	return revs
}

func AnimeReviewsWithUserInfoDomain(animeId int, userId string) []TReviewJoinUser {
	ctx := context.Background()
	rows := AnimeReviewsList(animeId, userId)
	var reviews []TReviewJoinUser
	for rows.Next() {
		var r TReviewJoinUser
		err := rows.Scan(
			&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId, &r.CreatedAt, &r.UpdatedAt,
		)
		if r.UserId != nil {
			user := auth.GetUserFirebase(ctx, *r.UserId)
			r.User = user
		} else {
			r.User = nil
		}

		if err != nil {
			continue
			//fmt.Print(err)
		}
		reviews = append(reviews, r)
	}
	defer rows.Close()
	return reviews
}

func OnesReviewsJoinAnimeDomain(userId string) []TReviewJoinAnime {
	rows := OnesReviewsJoinAnime(userId)
	var revs []TReviewJoinAnime
	for rows.Next() {
		var r TReviewJoinAnime
		err := rows.Scan(
			&r.ID, &r.Content, &r.Rating, &r.AnimeId, &r.UserId,
			&r.CreatedAt, &r.UpdatedAt, &r.Title, &r.Slug, &r.AnimeContent, &r.AState,
		)
		if err != nil {
			panic(err.Error())
		}
		revs = append(revs, r)
	}
	defer rows.Close()
	return revs
}
