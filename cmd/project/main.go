package main

import (
	"animar/v1/anime"
	"animar/v1/auth"
	"animar/v1/blog"
	"animar/v1/configs"
	"animar/v1/platform"
	"animar/v1/review"
	"animar/v1/seasons"
	"animar/v1/series"
	"animar/v1/tools"
	"animar/v1/watch"
	"fmt"
	"net/http"
)

func main() {

	configs.SetEnviron()
	defer configs.SetEnviron()

	// connection
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
		fmt.Print("OK")
	})

	/*   Anime database   */
	http.HandleFunc("/db/anime/", tools.Handle(anime.AnimeView))
	http.HandleFunc("/db/anime/search/", tools.Handle(anime.SearchAnimeView))
	http.HandleFunc("/db/anime/minimum/", tools.Handle(anime.ListAnimeMinimumView))
	http.HandleFunc("/db/anime/wc/", tools.Handle(anime.AnimeWithUserWatchView)) // 動かず

	/*   blogs   */
	http.HandleFunc("/blog/", tools.Handle(blog.BlogJoinAnimeView))
	http.HandleFunc("/blog/post/", tools.Handle(blog.InsertBlogWithRelationView))
	http.HandleFunc("/blog/delete/", tools.Handle(blog.DeleteBlogView))             // ?id=
	http.HandleFunc("/blog/update/", tools.Handle(blog.UpdateBlogWithRelationView)) // ?id=

	/*   reviews   */
	http.HandleFunc("/reviews/", tools.Handle(review.GetYourReviews))
	http.HandleFunc("/reviews/user/", tools.Handle(review.GetOnesReviewsView))
	http.HandleFunc("/reviews/post/star/", tools.Handle(review.UpsertReviewStarView))
	http.HandleFunc("/reviews/post/content/", tools.Handle(review.UpsertReviewContentView))
	//http.HandleFunc("/reviews/anime/", tools.Handle(review.GetAnimeReviews))
	http.HandleFunc("/reviews/anime/", tools.Handle(review.GetAnimeReviewsView))
	http.HandleFunc("/reviews/ua/", tools.Handle(review.GetAnimeUserReviewView))
	http.HandleFunc("/reviews/star/", tools.Handle(review.AnimeStarAvgView)) // star average

	/*   watches count   */
	http.HandleFunc("/watch/", tools.Handle(watch.AnimeWatchCountView))
	http.HandleFunc("/watch/u/", tools.Handle(watch.UserWatchStatusView))
	http.HandleFunc("/watch/post/", tools.Handle(watch.WatchPostView)) // upsert
	http.HandleFunc("/watch/delete/", tools.Handle(watch.WatchDeleteView))
	http.HandleFunc("/watch/ua/", tools.Handle(watch.WatchAnimeStateOfUserView))

	/*   auth   */
	http.HandleFunc("/auth/user/", tools.Handle(auth.GetUserModelView)) // ?uid=<UID>
	http.HandleFunc("/auth/login/post/", tools.Handle(auth.SetJWTCookieView))
	http.HandleFunc("/auth/refresh/", tools.Handle(auth.RenewTokenFCView))
	http.HandleFunc("/auth/cookie/", tools.Handle(auth.TestGetCookie))
	http.HandleFunc("/auth/user/cookie/", tools.Handle(auth.GetUserModelFCView))
	//http.HandleFunc("/auth/user/cookie/", tools.Handle(auth.GetUserModelFCWithVerifiedView))
	http.HandleFunc("/auth/register/", tools.Handle(auth.CreateUserFirstView))
	http.HandleFunc("/auth/profile/update/", tools.Handle(auth.UserUpdateView))

	/*   admin   */
	http.HandleFunc("/admin/anime/", tools.Handle(anime.AnimeListAdminView))
	http.HandleFunc("/admin/anime/detail/", tools.Handle(anime.AnimeDetailAdminView))
	http.HandleFunc("/admin/anime/post/", tools.Handle(anime.AnimePostView))
	http.HandleFunc("/admin/anime/update/", tools.Handle(anime.AnimeUpdateView))
	http.HandleFunc("/admin/anime/delete/", tools.Handle(anime.AnimeDeleteView))

	/*   series   */
	http.HandleFunc("/series/", tools.Handle(series.SeriesView))
	http.HandleFunc("/series/post/", tools.Handle(series.InsertSeriesView))

	/*   seasons   */
	http.HandleFunc("/season/", tools.Handle(seasons.SeasonView))
	http.HandleFunc("/admin/season/post/", tools.Handle(seasons.InsertSeasonView))
	// relations
	http.HandleFunc("/admin/season/anime/post/", tools.Handle(seasons.InsertRelationSeasonView))
	http.HandleFunc("/season/anime/", tools.Handle(seasons.SeasonByAnimeIdView)) // ?id=

	/*   platform   */
	http.HandleFunc("/admin/platform/", tools.Handle(platform.PlatformView))
	http.HandleFunc("/admin/platform/post/", tools.Handle(platform.InsertPlatformView))
	http.HandleFunc("/admin/platform/update/", tools.Handle(platform.UpdatePlatformView))
	http.HandleFunc("/admin/platform/delete/", tools.Handle(platform.DeletePlatformView))
	// relations
	http.HandleFunc("/relation/plat/", tools.Handle(platform.RelationPlatformView)) // ?id=<anime_id>
	http.HandleFunc("/admin/relation/plat/post/", tools.Handle(platform.InsertRelationPlatformView))
	http.HandleFunc("/admin/relation/plat/delete/", tools.Handle(platform.DeleteRelationPlatformView)) // ?anime=<anime_id>&platform=<platform_id>

	if tools.IsProductionEnv() {
		//http.ListenAndServeTLS(":443", os.Getenv("SSL_CHAIN_PATH"), os.Getenv("SSL_KEY_PATH"), nil)
		http.ListenAndServe(":8000", nil) // reverse proxy
	} else {
		http.ListenAndServe(":8000", nil)
	}

}
