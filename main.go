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
	http.HandleFunc("/db/anime/post/", helper.Handle(anime.AnimePostView))

	/*   reviews   */
	http.HandleFunc("/reviews/", helper.Handle(review.GetYourReviews))
	http.HandleFunc("/reviews/post/", helper.Handle(review.ReviewPostView))
	http.HandleFunc("/reviews/anime/", helper.Handle(review.GetAnimeReviews))

	/*   watches count   */
	http.HandleFunc("/watch/", helper.Handle(watch.AnimeWatchCountView))
	http.HandleFunc("/watch/u/", helper.Handle(watch.UserWatchStatusView))
	http.HandleFunc("/watch/post/", helper.Handle(watch.WatchPostView))
	http.HandleFunc("/watch/delete/", helper.Handle(watch.WatchDeleteView))

	/*   auth   */
	http.HandleFunc("/auth/sample/", helper.Handle(auth.SampleGetUserJson)) // ?uid=<UID>
	http.HandleFunc("/auth/login/post/", helper.Handle(auth.SetJWTCookie))
	http.HandleFunc("/auth/refresh/", helper.Handle(auth.RenewTokenView))
	http.HandleFunc("/auth/cookie/", helper.Handle(auth.TestGetCookie))

	http.ListenAndServe(":8000", nil)
}
