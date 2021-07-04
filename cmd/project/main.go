package main

import (
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/controllers"
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
	http.Handle("/db/anime/", animeController.BaseMiddleware(animeController.GiveUserIdMiddlewareAbleSSR(http.HandlerFunc(animeController.AnimeView))))
	http.Handle("/db/anime/search/", animeController.BaseMiddleware(http.HandlerFunc(animeController.SearchAnimeMinimumView)))
	http.Handle("/db/anime/minimum/", animeController.BaseMiddleware(http.HandlerFunc(animeController.AnimeMinimumsView)))

	/*   blogs   */
	blogController := controllers.NewBlogController(sqlHandler)
	http.Handle("/blog/", blogController.BaseMiddleware(blogController.GiveUserIdMiddlewareAbleSSR(http.HandlerFunc(blogController.BlogJoinAnimeView))))
	http.Handle("/blog/post/", blogController.BaseMiddleware(blogController.PostOnlyMiddleware(blogController.LoginRequireMiddleware(http.HandlerFunc(blogController.InsertBlogWithRelationView)))))
	http.Handle("/blog/delete/", blogController.BaseMiddleware(blogController.DeleteOnlyMiddleware(blogController.LoginRequireMiddleware(http.HandlerFunc(blogController.DeleteBlogView)))))          // ?id=
	http.Handle("/blog/update/", blogController.BaseMiddleware(blogController.PutOnlyMiddleware(blogController.LoginRequireMiddleware(http.HandlerFunc(blogController.UpdateBlogWithRelationView))))) // ?id=

	/*   reviews   */
	reviewController := controllers.NewReviewController(sqlHandler)
	http.Handle("/reviews/user/", reviewController.BaseMiddleware(http.HandlerFunc(reviewController.GetOnesReviewsView)))
	http.Handle("/reviews/post/star/", reviewController.BaseMiddleware(reviewController.UpsertOnlyMiddleware(reviewController.LoginRequireMiddleware(http.HandlerFunc(reviewController.UpsertReviewRatingView)))))
	http.Handle("/reviews/post/content/", reviewController.BaseMiddleware(reviewController.UpsertOnlyMiddleware(reviewController.LoginRequireMiddleware(http.HandlerFunc(reviewController.UpsertReviewContentView)))))
	http.Handle("/reviews/anime/", reviewController.BaseMiddleware(http.HandlerFunc(reviewController.GetAnimeReviewsView)))
	http.Handle("/reviews/ua/", reviewController.BaseMiddleware(reviewController.GiveUserIdMiddleware(http.HandlerFunc(reviewController.GetAnimeReviewOfUserView))))
	http.Handle("/reviews/star/", reviewController.BaseMiddleware(http.HandlerFunc(reviewController.AnimeRatingAvgView))) // star average

	/*   watches count   */
	audienceController := controllers.NewAudienceController(sqlHandler)
	http.Handle("/watch/", audienceController.BaseMiddleware(http.HandlerFunc(audienceController.AnimeAudienceCountsView)))
	http.Handle("/watch/u/", audienceController.BaseMiddleware(http.HandlerFunc(audienceController.AudienceWithAnimeByUserView)))
	http.Handle("/watch/post/", audienceController.BaseMiddleware(audienceController.UpsertOnlyMiddleware(audienceController.LoginRequireMiddleware(http.HandlerFunc(audienceController.UpsertAudienceView))))) // upsert
	http.Handle("/watch/delete/", audienceController.BaseMiddleware(audienceController.DeleteOnlyMiddleware(audienceController.LoginRequireMiddleware(http.HandlerFunc(audienceController.DeleteAudienceView)))))
	http.Handle("/watch/ua/", audienceController.BaseMiddleware(audienceController.GiveUserIdMiddleware(http.HandlerFunc(audienceController.AudienceByAnimeAndUserView))))

	/*   auth   */
	firebase := infrastructure.NewFireBaseClient()
	authController := controllers.NewAuthController(firebase)
	http.Handle("/auth/user/", authController.BaseMiddleware(http.HandlerFunc(authController.GetUserModelFromQueryView)))
	http.Handle("/auth/login/post/", authController.BaseMiddleware(authController.PostOnlyMiddleware(http.HandlerFunc(authController.LoginView))))
	http.Handle("/auth/refresh/", authController.BaseMiddleware(http.HandlerFunc(authController.RenewTokenView)))
	http.Handle("/auth/user/cookie/", authController.BaseMiddleware(authController.LoginRequireMiddleware(http.HandlerFunc(authController.GetUserModelFromCookieView))))
	http.Handle("/auth/register/", authController.BaseMiddleware(authController.PostOnlyMiddleware(http.HandlerFunc(authController.RegisterView))))
	http.Handle("/auth/profile/update/", authController.BaseMiddleware(authController.PostOnlyMiddleware(http.HandlerFunc(authController.UpdateProfileView))))
	// oauth
	http.Handle("/auth/google/", authController.BaseMiddleware(http.HandlerFunc(authController.GoogleOAuthView)))
	http.Handle("/auth/google/redirect/", authController.BaseMiddleware(http.HandlerFunc(authController.GoogleRedirectView)))

	/*   admin   */
	adminController := controllers.NewAdminController(sqlHandler)
	http.Handle("/admin/anime/", adminController.BaseMiddleware(adminController.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.AnimeListAdminView))))
	http.Handle("/admin/anime/detail/", adminController.BaseMiddleware(adminController.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.AnimeDetailAdminView))))
	http.Handle("/admin/anime/post/", adminController.BaseMiddleware(adminController.PostOnlyMiddleware(adminController.AdminRequiredMiddleware(http.HandlerFunc(adminController.AnimePostAdminView)))))
	http.Handle("/admin/anime/update/", adminController.BaseMiddleware(adminController.PutOnlyMiddleware(adminController.AdminRequiredMiddleware(http.HandlerFunc(adminController.AnimeUpdateView)))))
	http.Handle("/admin/anime/delete/", adminController.BaseMiddleware(adminController.DeleteOnlyMiddleware(adminController.AdminRequiredMiddleware(http.HandlerFunc(adminController.AnimeDeleteView)))))

	/*   series   */
	http.Handle("/series/", adminController.BaseMiddleware(adminController.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.SeriesView))))
	http.Handle("/series/post/", adminController.BaseMiddleware(adminController.PostOnlyMiddleware(adminController.AdminRequiredMiddleware(http.HandlerFunc(adminController.InsertSeriesView)))))

	/*   seasons   */
	seasonController := controllers.NewSeasonController(sqlHandler)
	http.Handle("/season/", adminController.BaseMiddleware(adminController.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.SeasonView))))
	http.Handle("/admin/season/post/", adminController.BaseMiddleware(adminController.AdminRequiredMiddleware(adminController.PostOnlyMiddleware(http.HandlerFunc(adminController.InsertSeasonView)))))
	// relations
	http.Handle("/admin/season/anime/post/", adminController.BaseMiddleware(adminController.AdminRequiredMiddleware(adminController.PostOnlyMiddleware(http.HandlerFunc(adminController.InsertRelationSeasonView)))))
	http.Handle("/season/anime/", seasonController.BaseMiddleware(http.HandlerFunc(seasonController.SeasonByAnimeIdView))) // ?id=

	/*   platform   */
	http.Handle("/admin/platform/", adminController.BaseMiddleware(adminController.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.PlatformView))))
	http.Handle("/admin/platform/post/", adminController.BaseMiddleware(adminController.PostOnlyMiddleware(adminController.AdminRequiredMiddleware(http.HandlerFunc(adminController.PlatformInsertView)))))
	http.Handle("/admin/platform/update/", adminController.BaseMiddleware(adminController.PutOnlyMiddleware(adminController.AdminRequiredMiddleware(http.HandlerFunc(adminController.PlatformUpdateView)))))
	http.Handle("/admin/platform/delete/", adminController.BaseMiddleware(adminController.DeleteOnlyMiddleware(adminController.AdminRequiredMiddleware(http.HandlerFunc(adminController.PlatformDeleteview)))))
	// relations
	http.Handle("/relation/plat/", adminController.BaseMiddleware(http.HandlerFunc(adminController.RelationPlatformView))) // ?id=<anime_id>
	http.Handle("/admin/relation/plat/post/", adminController.BaseMiddleware(adminController.PostOnlyMiddleware(adminController.AdminRequiredMiddleware(http.HandlerFunc(adminController.InsertRelationPlatformView)))))
	http.Handle("/admin/relation/plat/delete/", adminController.BaseMiddleware(adminController.DeleteOnlyMiddleware(adminController.AdminRequiredMiddleware(http.HandlerFunc(adminController.DeleteRelationPlatformView))))) // ?anime=<anime_id>&platform=<platform_id>

	/*   utilities   */
	utilityController := controllers.NewUtilityController()
	http.Handle("/utils/s3/", utilityController.BaseMiddleware(utilityController.PostOnlyMiddleware(utilityController.LoginRequireMiddleware(http.HandlerFunc(utilityController.SimpleUploadImage)))))

	if tools.IsProductionEnv() {
		http.ListenAndServe(":8000", nil) // reverse proxy
	} else {
		http.ListenAndServe(":8000", nil)
	}

}
