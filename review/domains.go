package review

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
