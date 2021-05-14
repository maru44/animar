package watch

func AnimeWatchCountDomain(animeId int) []TWatchCount {
	rows := AnimeWatchCounts(animeId)
	var watches []TWatchCount
	for rows.Next() {
		var watch TWatchCount
		err := rows.Scan(&watch.State, &watch.Count)
		if err != nil {
			panic(err.Error())
		}
		watches = append(watches, watch)
	}
	defer rows.Close()
	return watches
}

func OnesWatchStatusDomain(userId string) []TWatch {
	rows := OnesAnimeWatchList(userId)
	var watches []TWatch
	for rows.Next() {
		var w TWatch
		err := rows.Scan(
			&w.ID, &w.Watch, &w.AnimeId, &w.UserId, &w.CreatedAt, &w.UpdatedAt,
		)
		if err != nil {
			panic(err.Error())
		}
		watches = append(watches, w)
	}
	defer rows.Close()
	return watches
}
