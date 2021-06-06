package watch

func AnimeWatchCountDomain(animeId int) []TAudienceCount {
	rows := AnimeWatchCounts(animeId)
	var watches []TAudienceCount
	for rows.Next() {
		var w TAudienceCount
		err := rows.Scan(&w.State, &w.Count)
		if err != nil {
			panic(err.Error())
		}
		watches = append(watches, w)
	}
	defer rows.Close()
	return watches
}

func OnesWatchStatusDomain(userId string) []TAudienceJoinAnime {
	rows := OnesAnimeWatchJoinList(userId)
	var watches []TAudienceJoinAnime
	for rows.Next() {
		var w TAudienceJoinAnime
		err := rows.Scan(
			&w.ID, &w.State, &w.AnimeId, &w.UserId, &w.CreatedAt,
			&w.UpdatedAt, &w.Title, &w.Slug, &w.Content, &w.AState,
		)
		if err != nil {
			panic(err.Error())
		}
		watches = append(watches, w)
	}
	defer rows.Close()
	return watches
}
