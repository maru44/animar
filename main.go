package main

import (
	"animar/v1/anime"
	"animar/v1/auth"
	"animar/v1/helper"
	"animar/v1/review"
	"animar/v1/watch"
	"net/http"
)

func main() {

	helper.SetEnviron()

	/*   Anime database   */
	http.HandleFunc("/db/anime/", helper.Handle(anime.AnimeView))
	http.HandleFunc("/db/anime/wc/", helper.Handle(anime.AnimeWithUserWatchView)) // 動かず
	http.HandleFunc("/db/anime/post/", helper.Handle(anime.AnimePostView))

	/*   reviews   */
	http.HandleFunc("/reviews/", helper.Handle(review.GetYourReviews))
	http.HandleFunc("/reviews/post/", helper.Handle(review.ReviewPostView))
	http.HandleFunc("/reviews/post/star/", helper.Handle(review.UpsertReviewStarView))
	http.HandleFunc("/reviews/post/content/", helper.Handle(review.UpsertReviewContentView))
	http.HandleFunc("/reviews/anime/", helper.Handle(review.GetAnimeReviews))

	/*   watches count   */
	http.HandleFunc("/watch/", helper.Handle(watch.AnimeWatchCountView))
	http.HandleFunc("/watch/u/", helper.Handle(watch.UserWatchStatusView))
	http.HandleFunc("/watch/post/", helper.Handle(watch.WatchPostView)) // upsert
	http.HandleFunc("/watch/delete/", helper.Handle(watch.WatchDeleteView))
	http.HandleFunc("/watch/ua/", helper.Handle(watch.WatchAnimeStateOfUserView))

	/*   auth   */
	http.HandleFunc("/auth/sample/", helper.Handle(auth.SampleGetUserJsonView)) // ?uid=<UID>
	http.HandleFunc("/auth/login/post/", helper.Handle(auth.SetJWTCookieView))
	http.HandleFunc("/auth/refresh/", helper.Handle(auth.RenewTokenFCView))
	http.HandleFunc("/auth/cookie/", helper.Handle(auth.TestGetCookie))
	http.HandleFunc("/auth/user/cookie/", helper.Handle(auth.GetUserModelFCView))

	http.ListenAndServe(":8000", nil)
}
