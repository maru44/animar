package main

import (
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/controllers"
	"animar/v1/pkg/tools/tools"
	"net/http"
)

func main() {

	// connection
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	host, _ := os.Hostname()
	// 	w.Write([]byte(host))
	// })

	sqlHandler := infrastructure.NewSqlHandler()
	base := controllers.NewBaseController()

	http.Handle("/", base.BaseMiddleware(http.HandlerFunc(base.GatewayView)))

	/*   Anime database   */
	// http.HandleFunc("/db/anime/", handler.Handle(anime.AnimeView))
	animeController := controllers.NewAnimeController(sqlHandler)
	http.Handle("/db/anime/", base.BaseMiddleware(base.GiveUserIdMiddlewareAbleSSR(http.HandlerFunc(animeController.AnimeView))))
	http.Handle("/db/anime/search/", base.BaseMiddleware(http.HandlerFunc(animeController.SearchAnimeMinimumView)))
	http.Handle("/db/anime/minimum/", base.BaseMiddleware(http.HandlerFunc(animeController.AnimeMinimumsView)))

	/*   blogs   */
	blogController := controllers.NewBlogController(sqlHandler)
	http.Handle("/blog/", base.BaseMiddleware(base.GiveUserIdMiddlewareAbleSSR(http.HandlerFunc(blogController.BlogJoinAnimeView))))
	http.Handle("/blog/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(blogController.InsertBlogWithRelationView)))))
	http.Handle("/blog/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(blogController.DeleteBlogView)))))          // ?id=
	http.Handle("/blog/update/", base.BaseMiddleware(base.PutOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(blogController.UpdateBlogWithRelationView))))) // ?id=

	/*   reviews   */
	reviewController := controllers.NewReviewController(sqlHandler)
	http.Handle("/reviews/user/", base.BaseMiddleware(http.HandlerFunc(reviewController.GetOnesReviewsView)))
	http.Handle("/reviews/post/star/", base.BaseMiddleware(base.UpsertOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(reviewController.UpsertReviewRatingView)))))
	http.Handle("/reviews/post/content/", base.BaseMiddleware(base.UpsertOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(reviewController.UpsertReviewContentView)))))
	http.Handle("/reviews/anime/", base.BaseMiddleware(http.HandlerFunc(reviewController.GetAnimeReviewsView)))
	http.Handle("/reviews/ua/", base.BaseMiddleware(base.GiveUserIdMiddleware(http.HandlerFunc(reviewController.GetAnimeReviewOfUserView))))
	http.Handle("/reviews/star/", base.BaseMiddleware(http.HandlerFunc(reviewController.AnimeRatingAvgView))) // star average

	/*   watches count   */
	audienceController := controllers.NewAudienceController(sqlHandler)
	http.Handle("/watch/", base.BaseMiddleware(http.HandlerFunc(audienceController.AnimeAudienceCountsView)))
	http.Handle("/watch/u/", base.BaseMiddleware(http.HandlerFunc(audienceController.AudienceWithAnimeByUserView)))
	http.Handle("/watch/post/", base.BaseMiddleware(base.UpsertOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(audienceController.UpsertAudienceView))))) // upsert
	http.Handle("/watch/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(audienceController.DeleteAudienceView)))))
	http.Handle("/watch/ua/", base.BaseMiddleware(base.GiveUserIdMiddleware(http.HandlerFunc(audienceController.AudienceByAnimeAndUserView))))

	/*   auth   */
	firebase := infrastructure.NewFireBaseClient()
	authController := controllers.NewAuthController(firebase)
	http.Handle("/auth/user/", base.BaseMiddleware(http.HandlerFunc(authController.GetUserModelFromQueryView)))
	http.Handle("/auth/login/post/", base.BaseMiddleware(base.PostOnlyMiddleware(http.HandlerFunc(authController.LoginView))))
	http.Handle("/auth/refresh/", base.BaseMiddleware(http.HandlerFunc(authController.RenewTokenView)))
	http.Handle("/auth/user/cookie/", base.BaseMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(authController.GetUserModelFromCookieView))))
	http.Handle("/auth/register/", base.BaseMiddleware(base.PostOnlyMiddleware(http.HandlerFunc(authController.RegisterView))))
	http.Handle("/auth/profile/update/", base.BaseMiddleware(base.PostOnlyMiddleware(base.GiveUserIdMiddleware(http.HandlerFunc(authController.UpdateProfileView)))))
	// oauth
	http.Handle("/auth/setcookie/", base.BaseMiddleware(base.PostOnlyMiddleware(http.HandlerFunc(authController.SetJwtTokenView))))

	/*   admin   */
	adminController := controllers.NewAdminController(sqlHandler)
	http.Handle("/admin/anime/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.AnimeListAdminView))))
	http.Handle("/admin/anime/detail/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.AnimeDetailAdminView))))
	http.Handle("/admin/anime/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.AnimePostAdminView)))))
	http.Handle("/admin/anime/update/", base.BaseMiddleware(base.PutOnlyMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.AnimeUpdateView)))))
	http.Handle("/admin/anime/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.AnimeDeleteView)))))

	/*   series   */
	http.Handle("/series/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.SeriesView))))
	http.Handle("/series/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.InsertSeriesView)))))

	/*   seasons   */
	seasonController := controllers.NewSeasonController(sqlHandler)
	http.Handle("/season/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.SeasonView))))
	http.Handle("/admin/season/post/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.PostOnlyMiddleware(http.HandlerFunc(adminController.InsertSeasonView)))))
	// relations
	http.Handle("/admin/season/anime/post/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.PostOnlyMiddleware(http.HandlerFunc(adminController.InsertRelationSeasonView)))))
	http.Handle("/season/anime/", base.BaseMiddleware(http.HandlerFunc(seasonController.SeasonByAnimeIdView))) // ?id=

	/*   platform   */
	http.Handle("/admin/platform/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.PlatformView))))
	http.Handle("/admin/platform/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.PlatformInsertView)))))
	http.Handle("/admin/platform/update/", base.BaseMiddleware(base.PutOnlyMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.PlatformUpdateView)))))
	http.Handle("/admin/platform/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.PlatformDeleteview)))))
	// relations
	http.Handle("/relation/plat/", base.BaseMiddleware(http.HandlerFunc(adminController.RelationPlatformView))) // ?id=<anime_id>
	http.Handle("/admin/relation/plat/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.InsertRelationPlatformView)))))
	http.Handle("/admin/relation/plat/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.DeleteRelationPlatformView))))) // ?anime=<anime_id>&platform=<platform_id>

	/*   utilities   */
	utilityController := controllers.NewUtilityController()
	http.Handle("/utils/s3/", base.BaseMiddleware(base.PostOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(utilityController.SimpleUploadImage)))))

	if tools.IsProductionEnv() {
		http.ListenAndServe(":8000", nil) // reverse proxy
	} else {
		http.ListenAndServe(":8000", nil)
	}

}
