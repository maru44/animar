package main

import (
	"animar/v1/anime"
	"animar/v1/auth"
	"animar/v1/blog"
	"animar/v1/platform"
	"animar/v1/review"
	"animar/v1/seasons"
	"animar/v1/series"
	"animar/v1/tools"
	"animar/v1/watch"
	"net/http"
	"os"
)

func main() {

	// connection
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host, _ := os.Hostname()

		w.Write([]byte(host))
	})

	/*   Anime database   */
	http.HandleFunc("/db/anime/", tools.Handle(anime.AnimeView))
	http.HandleFunc("/db/anime/search/", tools.Handle(anime.SearchAnimeView))
	http.HandleFunc("/db/anime/minimum/", tools.Handle(anime.ListAnimeMinimumView))
	http.HandleFunc("/db/anime/wc/", tools.Handle(anime.AnimeWithUserWatchView)) // 動かず

	/*   blogs   */
	http.HandleFunc("/blog/", tools.Handle(blog.BlogJoinAnimeView))
	http.HandleFunc("/blog/post/", tools.Handle(tools.PostOnlyMiddleware, blog.InsertBlogWithRelationView))
	http.HandleFunc("/blog/delete/", tools.Handle(tools.DeleteOnlyMiddleware, blog.DeleteBlogView))          // ?id=
	http.HandleFunc("/blog/update/", tools.Handle(tools.PutOnlyMiddleware, blog.UpdateBlogWithRelationView)) // ?id=

	/*   reviews   */
	http.HandleFunc("/reviews/", tools.Handle(review.GetYourReviews))
	http.HandleFunc("/reviews/user/", tools.Handle(review.GetOnesReviewsView))
	http.HandleFunc("/reviews/post/star/", tools.Handle(tools.UpsertOnlyMiddleware, review.UpsertReviewStarView))
	http.HandleFunc("/reviews/post/content/", tools.Handle(tools.UpsertOnlyMiddleware, review.UpsertReviewContentView))
	//http.HandleFunc("/reviews/anime/", tools.Handle(review.GetAnimeReviews))
	http.HandleFunc("/reviews/anime/", tools.Handle(review.GetAnimeReviewsView))
	http.HandleFunc("/reviews/ua/", tools.Handle(review.GetAnimeUserReviewView))
	http.HandleFunc("/reviews/star/", tools.Handle(review.AnimeStarAvgView)) // star average

	/*   watches count   */
	http.HandleFunc("/watch/", tools.Handle(watch.AnimeWatchCountView))
	http.HandleFunc("/watch/u/", tools.Handle(watch.UserWatchStatusView))
	http.HandleFunc("/watch/post/", tools.Handle(tools.PostOnlyMiddleware, watch.WatchPostView)) // upsert
	http.HandleFunc("/watch/delete/", tools.Handle(tools.DeleteOnlyMiddleware, watch.WatchDeleteView))
	http.HandleFunc("/watch/ua/", tools.Handle(watch.WatchAnimeStateOfUserView))

	/*   auth   */
	http.HandleFunc("/auth/user/", tools.Handle(auth.GetUserModelView)) // ?uid=<UID>
	http.HandleFunc("/auth/login/post/", tools.Handle(tools.PostOnlyMiddleware, auth.SetJWTCookieView))
	http.HandleFunc("/auth/refresh/", tools.Handle(auth.RenewTokenFCView))
	http.HandleFunc("/auth/cookie/", tools.Handle(auth.TestGetCookie))
	http.HandleFunc("/auth/user/cookie/", tools.Handle(auth.GetUserModelFCView))
	//http.HandleFunc("/auth/user/cookie/", tools.Handle(auth.GetUserModelFCWithVerifiedView))
	http.HandleFunc("/auth/register/", tools.Handle(tools.PostOnlyMiddleware, auth.CreateUserFirstView))
	http.HandleFunc("/auth/profile/update/", tools.Handle(tools.PostOnlyMiddleware, auth.UserUpdateView))

	/*   admin   */
	http.HandleFunc("/admin/anime/", tools.Handle(tools.AdminRequiredMiddlewareGet, anime.AnimeListAdminView))
	http.HandleFunc("/admin/anime/detail/", tools.Handle(tools.AdminRequiredMiddlewareGet, anime.AnimeDetailAdminView))
	http.HandleFunc("/admin/anime/post/", tools.Handle(tools.PostOnlyMiddleware, tools.AdminRequiredMiddleware, anime.AnimePostView))
	http.HandleFunc("/admin/anime/update/", tools.Handle(tools.PutOnlyMiddleware, tools.AdminRequiredMiddleware, anime.AnimeUpdateView))
	http.HandleFunc("/admin/anime/delete/", tools.Handle(tools.DeleteOnlyMiddleware, tools.AdminRequiredMiddleware, anime.AnimeDeleteView))

	/*   series   */
	http.HandleFunc("/series/", tools.Handle(series.SeriesView))
	http.HandleFunc("/series/post/", tools.Handle(tools.PostOnlyMiddleware, tools.AdminRequiredMiddleware, series.InsertSeriesView))

	/*   seasons   */
	http.HandleFunc("/season/", tools.Handle(seasons.SeasonView))
	http.HandleFunc("/admin/season/post/", tools.Handle(tools.PostOnlyMiddleware, tools.PostOnlyMiddleware, seasons.InsertSeasonView))
	// relations
	http.HandleFunc("/admin/season/anime/post/", tools.Handle(tools.PostOnlyMiddleware, tools.AdminRequiredMiddleware, seasons.InsertRelationSeasonView))
	http.HandleFunc("/season/anime/", tools.Handle(seasons.SeasonByAnimeIdView)) // ?id=

	/*   platform   */
	http.HandleFunc("/admin/platform/", tools.Handle(tools.AdminRequiredMiddlewareGet, platform.PlatformView))
	http.HandleFunc("/admin/platform/post/", tools.Handle(tools.PostOnlyMiddleware, tools.AdminRequiredMiddleware, platform.InsertPlatformView))
	http.HandleFunc("/admin/platform/update/", tools.Handle(tools.PutOnlyMiddleware, tools.AdminRequiredMiddleware, platform.UpdatePlatformView))
	http.HandleFunc("/admin/platform/delete/", tools.Handle(tools.DeleteOnlyMiddleware, tools.AdminRequiredMiddleware, platform.DeletePlatformView))
	// relations
	http.HandleFunc("/relation/plat/", tools.Handle(platform.RelationPlatformView)) // ?id=<anime_id>
	http.HandleFunc("/admin/relation/plat/post/", tools.Handle(tools.PostOnlyMiddleware, tools.AdminRequiredMiddleware, platform.InsertRelationPlatformView))
	http.HandleFunc("/admin/relation/plat/delete/", tools.Handle(tools.DeleteOnlyMiddleware, tools.AdminRequiredMiddleware, platform.DeleteRelationPlatformView)) // ?anime=<anime_id>&platform=<platform_id>

	if tools.IsProductionEnv() {
		//http.ListenAndServeTLS(":443", os.Getenv("SSL_CHAIN_PATH"), os.Getenv("SSL_KEY_PATH"), nil)
		http.ListenAndServe(":8000", nil) // reverse proxy
	} else {
		http.ListenAndServe(":8000", nil)
	}

}
