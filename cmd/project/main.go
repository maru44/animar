package main

import (
	"animar/v1/pkg/controllers"
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/mvc/anime"
	"animar/v1/pkg/mvc/auth"
	"animar/v1/pkg/mvc/blog"
	"animar/v1/pkg/mvc/platform"
	"animar/v1/pkg/mvc/review"
	"animar/v1/pkg/mvc/seasons"
	"animar/v1/pkg/mvc/series"
	"animar/v1/pkg/mvc/watch"
	"animar/v1/pkg/tools/handler"
	"animar/v1/pkg/tools/middleware"
	"animar/v1/pkg/tools/tools"
	"net/http"
	"os"
)

func main() {

	// connection
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host, _ := os.Hostname()
		w.Write([]byte(host))
	})

	sqlHandler := infrastructure.NewSqlHandler()
	animeController := controllers.NewAnimeController(sqlHandler)

	http.HandleFunc("/db/anime/", handler.Handle(animeController.AnimeView))

	/*   Anime database   */
	// http.HandleFunc("/db/anime/", handler.Handle(anime.AnimeView))
	http.HandleFunc("/db/anime/search/", handler.Handle(anime.SearchAnimeMinView))
	http.HandleFunc("/db/anime/minimum/", handler.Handle(anime.ListAnimeMinimumView))

	/*   blogs   */
	http.HandleFunc("/blog/", handler.Handle(blog.BlogJoinAnimeView))
	http.HandleFunc("/blog/post/", handler.Handle(middleware.PostOnlyMiddleware, blog.InsertBlogWithRelationView))
	http.HandleFunc("/blog/delete/", handler.Handle(middleware.DeleteOnlyMiddleware, blog.DeleteBlogView))          // ?id=
	http.HandleFunc("/blog/update/", handler.Handle(middleware.PutOnlyMiddleware, blog.UpdateBlogWithRelationView)) // ?id=

	/*   reviews   */
	http.HandleFunc("/reviews/", handler.Handle(review.GetYourReviews))
	http.HandleFunc("/reviews/user/", handler.Handle(review.GetOnesReviewsView))
	http.HandleFunc("/reviews/post/star/", handler.Handle(middleware.UpsertOnlyMiddleware, review.UpsertReviewStarView))
	http.HandleFunc("/reviews/post/content/", handler.Handle(middleware.UpsertOnlyMiddleware, review.UpsertReviewContentView))
	//http.HandleFunc("/reviews/anime/", handler.Handle(review.GetAnimeReviews))
	http.HandleFunc("/reviews/anime/", handler.Handle(review.GetAnimeReviewsView))
	http.HandleFunc("/reviews/ua/", handler.Handle(review.GetAnimeUserReviewView))
	http.HandleFunc("/reviews/star/", handler.Handle(review.AnimeStarAvgView)) // star average

	/*   watches count   */
	http.HandleFunc("/watch/", handler.Handle(watch.AnimeWatchCountView))
	http.HandleFunc("/watch/u/", handler.Handle(watch.UserWatchStatusView))
	http.HandleFunc("/watch/post/", handler.Handle(middleware.PostOnlyMiddleware, watch.WatchPostView)) // upsert
	http.HandleFunc("/watch/delete/", handler.Handle(middleware.DeleteOnlyMiddleware, watch.WatchDeleteView))
	http.HandleFunc("/watch/ua/", handler.Handle(watch.WatchAnimeStateOfUserView))

	/*   auth   */
	http.HandleFunc("/auth/user/", handler.Handle(auth.GetUserModelView)) // ?uid=<UID>
	http.HandleFunc("/auth/login/post/", handler.Handle(middleware.PostOnlyMiddleware, auth.SetJWTCookieView))
	http.HandleFunc("/auth/refresh/", handler.Handle(auth.RenewTokenFCView))
	http.HandleFunc("/auth/cookie/", handler.Handle(auth.TestGetCookie))
	http.HandleFunc("/auth/user/cookie/", handler.Handle(auth.GetUserModelFCView))
	//http.HandleFunc("/auth/user/cookie/", handler.Handle(auth.GetUserModelFCWithVerifiedView))
	http.HandleFunc("/auth/register/", handler.Handle(middleware.PostOnlyMiddleware, auth.CreateUserFirstView))
	http.HandleFunc("/auth/profile/update/", handler.Handle(middleware.PostOnlyMiddleware, auth.UserUpdateView))

	/*   admin   */
	http.HandleFunc("/admin/anime/", handler.Handle(middleware.AdminRequiredMiddlewareGet, anime.AnimeListAdminView))
	http.HandleFunc("/admin/anime/detail/", handler.Handle(middleware.AdminRequiredMiddlewareGet, anime.AnimeDetailAdminView))
	http.HandleFunc("/admin/anime/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.AdminRequiredMiddleware, anime.AnimePostView))
	http.HandleFunc("/admin/anime/update/", handler.Handle(middleware.PutOnlyMiddleware, middleware.AdminRequiredMiddleware, anime.AnimeUpdateView))
	http.HandleFunc("/admin/anime/delete/", handler.Handle(middleware.DeleteOnlyMiddleware, middleware.AdminRequiredMiddleware, anime.AnimeDeleteView))

	/*   series   */
	http.HandleFunc("/series/", handler.Handle(middleware.AdminRequiredMiddlewareGet, series.SeriesView))
	http.HandleFunc("/series/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.AdminRequiredMiddleware, series.InsertSeriesView))

	/*   seasons   */
	http.HandleFunc("/season/", handler.Handle(middleware.AdminRequiredMiddlewareGet, seasons.SeasonView))
	http.HandleFunc("/admin/season/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.PostOnlyMiddleware, seasons.InsertSeasonView))
	// relations
	http.HandleFunc("/admin/season/anime/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.AdminRequiredMiddleware, seasons.InsertRelationSeasonView))
	http.HandleFunc("/season/anime/", handler.Handle(seasons.SeasonByAnimeIdView)) // ?id=

	/*   platform   */
	http.HandleFunc("/admin/platform/", handler.Handle(middleware.AdminRequiredMiddlewareGet, platform.PlatformView))
	http.HandleFunc("/admin/platform/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.AdminRequiredMiddleware, platform.InsertPlatformView))
	http.HandleFunc("/admin/platform/update/", handler.Handle(middleware.PutOnlyMiddleware, middleware.AdminRequiredMiddleware, platform.UpdatePlatformView))
	http.HandleFunc("/admin/platform/delete/", handler.Handle(middleware.DeleteOnlyMiddleware, middleware.AdminRequiredMiddleware, platform.DeletePlatformView))
	// relations
	http.HandleFunc("/relation/plat/", handler.Handle(platform.RelationPlatformView)) // ?id=<anime_id>
	http.HandleFunc("/admin/relation/plat/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.AdminRequiredMiddleware, platform.InsertRelationPlatformView))
	http.HandleFunc("/admin/relation/plat/delete/", handler.Handle(middleware.DeleteOnlyMiddleware, middleware.AdminRequiredMiddleware, platform.DeleteRelationPlatformView)) // ?anime=<anime_id>&platform=<platform_id>

	if tools.IsProductionEnv() {
		//http.ListenAndServeTLS(":443", os.Getenv("SSL_CHAIN_PATH"), os.Getenv("SSL_KEY_PATH"), nil)
		http.ListenAndServe(":8000", nil) // reverse proxy
	} else {
		http.ListenAndServe(":8000", nil)
	}

}
