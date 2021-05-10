package main

import (
	"animar/v1/anime"
	"animar/v1/auth"
	"animar/v1/helper"
	"animar/v1/review"
	"net/http"
)

func main() {

	helper.SetEnviron()

	/*   Anime database   */
	http.HandleFunc("/db/anime/", helper.Handle(anime.AnimeView))
	http.HandleFunc("/db/anime/post/", helper.Handle(anime.AnimePostView))

	/*   reviews   */
	http.HandleFunc("/reviews/post/", helper.Handle(review.ReviewPostView))
	http.HandleFunc("/reviews/", helper.Handle(review.GetYourReviews))

	/*   auth   */
	//http.HandleFunc("/auth/sample/", auth.SampleGetUser)
	http.HandleFunc("/auth/sample/", helper.Handle(auth.SampleGetUserJson)) // ?uid=<UID>
	http.HandleFunc("/auth/login/post/", helper.Handle(auth.SetJWTCookie))
	http.HandleFunc("/auth/refresh/", helper.Handle(auth.RenewTokenView))
	http.HandleFunc("/auth/cookie/", helper.Handle(auth.TestGetCookie))

	http.ListenAndServe(":8080", nil)
}
