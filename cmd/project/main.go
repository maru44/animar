package main

import (
	"animar/v1/pkg/controllers"
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/mvc/anime"
	"animar/v1/pkg/mvc/auth"
	"animar/v1/pkg/mvc/platform"
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

	/*   Anime database   */
	// http.HandleFunc("/db/anime/", handler.Handle(anime.AnimeView))
	animeController := controllers.NewAnimeController(sqlHandler)
	http.HandleFunc("/db/anime/", handler.Handle(animeController.AnimeView))
	http.HandleFunc("/db/anime/search/", handler.Handle(anime.SearchAnimeMinView))
	http.HandleFunc("/db/anime/minimum/", handler.Handle(anime.ListAnimeMinimumView))

	/*   blogs   */
	blogController := controllers.NewBlogController(sqlHandler)
	http.HandleFunc("/blog/", handler.Handle(blogController.BlogJoinAnimeView))
	http.HandleFunc("/blog/post/", handler.Handle(middleware.PostOnlyMiddleware, blogController.InsertBlogWithRelationView))
	http.HandleFunc("/blog/delete/", handler.Handle(middleware.DeleteOnlyMiddleware, blogController.DeleteBlogView))          // ?id=
	http.HandleFunc("/blog/update/", handler.Handle(middleware.PutOnlyMiddleware, blogController.UpdateBlogWithRelationView)) // ?id=

	/*   reviews   */
	reviewController := controllers.NewReviewController(sqlHandler)
	//http.HandleFunc("/reviews/", handler.Handle(reviewController.))
	http.HandleFunc("/reviews/user/", handler.Handle(reviewController.GetOnesReviewsView))
	http.HandleFunc("/reviews/post/star/", handler.Handle(middleware.UpsertOnlyMiddleware, reviewController.UpsertReviewRatingView))
	http.HandleFunc("/reviews/post/content/", handler.Handle(middleware.UpsertOnlyMiddleware, reviewController.UpsertReviewContentView))
	//http.HandleFunc("/reviews/anime/", handler.Handle(review.GetAnimeReviews))
	http.HandleFunc("/reviews/anime/", handler.Handle(reviewController.GetAnimeReviewsView))
	http.HandleFunc("/reviews/ua/", handler.Handle(reviewController.GetAnimeReviewOfUserView))
	http.HandleFunc("/reviews/star/", handler.Handle(reviewController.AnimeRatingAvgView)) // star average

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
	adminController := controllers.NewAdminController(sqlHandler)
	http.HandleFunc("/admin/anime/", handler.Handle(middleware.AdminRequiredMiddlewareGet, adminController.AnimeListAdminView))
	http.HandleFunc("/admin/anime/detail/", handler.Handle(middleware.AdminRequiredMiddlewareGet, adminController.AnimeDetailAdminView))
	http.HandleFunc("/admin/anime/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.AdminRequiredMiddleware, adminController.AnimePostAdminView))
	http.HandleFunc("/admin/anime/update/", handler.Handle(middleware.PutOnlyMiddleware, middleware.AdminRequiredMiddleware, adminController.AnimeUpdateView))
	http.HandleFunc("/admin/anime/delete/", handler.Handle(middleware.DeleteOnlyMiddleware, middleware.AdminRequiredMiddleware, adminController.AnimeDeleteView))

	/*   series   */
	http.HandleFunc("/series/", handler.Handle(middleware.AdminRequiredMiddlewareGet, adminController.SeriesView))
	http.HandleFunc("/series/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.AdminRequiredMiddleware, adminController.InsertSeriesView))

	/*   seasons   */
	seasonController := controllers.NewSeasonController(sqlHandler)
	http.HandleFunc("/season/", handler.Handle(middleware.AdminRequiredMiddlewareGet, adminController.SeasonView))
	http.HandleFunc("/admin/season/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.PostOnlyMiddleware, adminController.InsertSeasonView))
	// relations
	http.HandleFunc("/admin/season/anime/post/", handler.Handle(middleware.PostOnlyMiddleware, middleware.AdminRequiredMiddleware, adminController.InsertRelationSeasonView))
	http.HandleFunc("/season/anime/", handler.Handle(seasonController.SeasonByAnimeIdView)) // ?id=

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
