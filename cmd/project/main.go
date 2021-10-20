package main

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/infrastructure"
	"animar/v1/internal/pkg/interfaces/controllers"
	"animar/v1/internal/pkg/tools/tools"
	"net/http"
	"time"
)

func main() {

	router := http.NewServeMux()
	lg := domain.NewAccessLog()
	cache := domain.NewCahce()

	sqlHandler := infrastructure.NewSqlHandler()
	uploader := infrastructure.NewS3Uploader()
	base := controllers.NewBaseController(*cache)

	go func() {
		cache.DeleteRegularly(domain.CacheTypeCsrf, time.Minute)
	}()

	router.Handle("/", base.BaseMiddleware(http.HandlerFunc(base.GatewayView)))

	/*   Anime database   */
	// http.HandleFunc("/db/anime/", handler.Handle(anime.AnimeView))
	animeController := controllers.NewAnimeController(sqlHandler)
	router.Handle("/db/anime/", base.BaseMiddleware(base.GiveUserIdMiddlewareAbleSSR(http.HandlerFunc(animeController.AnimeView))))
	router.Handle("/db/anime/search/", base.BaseMiddleware(http.HandlerFunc(animeController.SearchAnimeMinimumView)))
	router.Handle("/db/anime/minimum/", base.BaseMiddleware(http.HandlerFunc(animeController.AnimeMinimumsView)))

	/*   blogs   */
	blogController := controllers.NewBlogController(sqlHandler, uploader)
	router.Handle("/blog/", base.BaseMiddleware(base.GiveUserIdMiddlewareAbleSSR(http.HandlerFunc(blogController.BlogJoinAnimeView))))
	router.Handle("/blog/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.LoginRequireMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(blogController.InsertBlogWithRelationView))))))
	router.Handle("/blog/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.LoginRequireMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(blogController.DeleteBlogView))))))          // ?id=
	router.Handle("/blog/update/", base.BaseMiddleware(base.PutOnlyMiddleware(base.LoginRequireMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(blogController.UpdateBlogWithRelationView)))))) // ?id=
	router.Handle("/blog/image/", base.BaseMiddleware(base.PostOnlyMiddleware(base.LoginRequireMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(blogController.SimpleUploadImage))))))

	/*   reviews   */
	reviewController := controllers.NewReviewController(sqlHandler)
	router.Handle("/reviews/", base.BaseMiddleware(http.HandlerFunc(reviewController.GetReviewView)))
	router.Handle("/reviews/user/", base.BaseMiddleware(http.HandlerFunc(reviewController.GetOnesReviewsView)))
	router.Handle("/reviews/post/star/", base.BaseMiddleware(base.UpsertOnlyMiddleware(base.LoginRequireMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(reviewController.UpsertReviewRatingView))))))
	router.Handle("/reviews/post/content/", base.BaseMiddleware(base.UpsertOnlyMiddleware(base.LoginRequireMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(reviewController.UpsertReviewContentView))))))
	router.Handle("/reviews/anime/", base.BaseMiddleware(http.HandlerFunc(reviewController.GetAnimeReviewsView)))
	router.Handle("/reviews/ua/", base.BaseMiddleware(base.GiveUserIdMiddleware(http.HandlerFunc(reviewController.GetAnimeReviewOfUserView))))
	router.Handle("/reviews/star/", base.BaseMiddleware(http.HandlerFunc(reviewController.AnimeRatingAvgView))) // star average

	/*   watches count   */
	audienceController := controllers.NewAudienceController(sqlHandler)
	router.Handle("/watch/", base.BaseMiddleware(http.HandlerFunc(audienceController.AnimeAudienceCountsView)))
	router.Handle("/watch/u/", base.BaseMiddleware(http.HandlerFunc(audienceController.AudienceWithAnimeByUserView)))
	router.Handle("/watch/post/", base.BaseMiddleware(base.UpsertOnlyMiddleware(base.LoginRequireMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(audienceController.UpsertAudienceView)))))) // upsert
	router.Handle("/watch/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.LoginRequireMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(audienceController.DeleteAudienceView))))))
	router.Handle("/watch/ua/", base.BaseMiddleware(base.GiveUserIdMiddleware(http.HandlerFunc(audienceController.AudienceByAnimeAndUserView))))

	staffController := controllers.NewStaffController(sqlHandler)
	router.Handle("/staff/", base.BaseMiddleware(http.HandlerFunc(staffController.StaffListView)))

	roleController := controllers.NewRoleController(sqlHandler)
	router.Handle("/staffrole/", base.BaseMiddleware(http.HandlerFunc(roleController.ListStaffRoleView))) // ?anime=

	platController := controllers.NewPlatformController(sqlHandler)
	router.Handle("/notification/register/", base.BaseMiddleware(base.PostOnlyMiddleware(base.LoginRequireMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(platController.RegisterNotifiedTargetView))))))
	router.Handle("/notification/user/slack/", base.BaseMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(platController.GetUsersChannelView))))

	// article
	articleController := controllers.NewArticleController(sqlHandler)
	router.Handle("/article/", base.BaseMiddleware(http.HandlerFunc(articleController.ArticleListView)))
	router.Handle("/article/detail/", base.BaseMiddleware(http.HandlerFunc(articleController.ArticleDetailView)))
	router.Handle("/article/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(articleController.InsertArticleView)))))
	router.Handle("/article/update/", base.BaseMiddleware(base.PutOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(articleController.UpdateArticleView)))))
	router.Handle("/article/character/", base.BaseMiddleware(http.HandlerFunc(articleController.ArticleCharacterView)))
	router.Handle("/article/character/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(articleController.InsertArticleCharaView)))))
	router.Handle("/article/character/update/", base.BaseMiddleware(base.PutOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(articleController.UpdateArticleCharaView)))))
	router.Handle("/article/character/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(articleController.DeleteArticleCharaView)))))
	router.Handle("/article/interview/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(articleController.InsertInterviewView)))))

	/*   auth   */
	firebase := infrastructure.NewFireBaseClient()
	authController := controllers.NewAuthController(firebase, uploader)
	router.Handle("/auth/user/", base.BaseMiddleware(http.HandlerFunc(authController.GetUserModelFromQueryView)))
	router.Handle("/auth/login/post/", base.BaseMiddleware(base.PostOnlyMiddleware(http.HandlerFunc(authController.LoginView))))
	router.Handle("/auth/refresh/", base.BaseMiddleware(http.HandlerFunc(authController.RenewTokenView)))
	router.Handle("/auth/user/cookie/", base.BaseMiddleware(base.LoginRequireMiddleware(http.HandlerFunc(authController.GetUserModelFromCookieView))))
	router.Handle("/auth/register/", base.BaseMiddleware(base.PostOnlyMiddleware(http.HandlerFunc(authController.RegisterView))))
	router.Handle("/auth/profile/update/", base.BaseMiddleware(base.PostOnlyMiddleware(base.GiveUserIdMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(authController.UpdateProfileView))))))
	// oauth
	router.Handle("/auth/setcookie/", base.BaseMiddleware(base.PostOnlyMiddleware(http.HandlerFunc(authController.SetJwtTokenView))))

	/*************************************************
	*                                                *
	*                     Admin                      *
	*                                                *
	*************************************************/

	adminController := controllers.NewAdminController(sqlHandler, uploader)
	router.Handle("/admin/anime/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.AnimeListAdminView))))
	router.Handle("/admin/anime/detail/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.AnimeDetailAdminView))))
	router.Handle("/admin/anime/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.AnimePostAdminView))))))
	router.Handle("/admin/anime/update/", base.BaseMiddleware(base.PutOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.AnimeUpdateView))))))
	router.Handle("/admin/anime/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.AnimeDeleteView))))))

	/*   series   */
	router.Handle("/series/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.SeriesView))))
	router.Handle("/series/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.InsertSeriesView))))))

	/*   seasons   */
	seasonController := controllers.NewSeasonController(sqlHandler)
	router.Handle("/season/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.SeasonView))))
	router.Handle("/admin/season/post/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.PostOnlyMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.InsertSeasonView))))))
	// relations
	router.Handle("/admin/season/anime/post/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.PostOnlyMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.InsertRelationSeasonView))))))
	router.Handle("/season/anime/", base.BaseMiddleware(http.HandlerFunc(seasonController.SeasonByAnimeIdView))) // ?id=
	router.Handle("/admin/season/anime/delete/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.DeleteOnlyMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.DeleteRelationSeasonView))))))

	/*   platform   */
	router.Handle("/admin/platform/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.PlatformView))))
	router.Handle("/admin/platform/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.PlatformInsertView))))))
	router.Handle("/admin/platform/update/", base.BaseMiddleware(base.PutOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.PlatformUpdateView))))))
	router.Handle("/admin/platform/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.PlatformDeleteview))))))
	// relations
	router.Handle("/relation/plat/", base.BaseMiddleware(http.HandlerFunc(adminController.RelationPlatformView))) // ?id=<anime_id>
	router.Handle("/admin/relation/plat/post/", base.BaseMiddleware(base.PostOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.InsertRelationPlatformView))))))
	router.Handle("/admin/relation/plat/delete/", base.BaseMiddleware(base.DeleteOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.DeleteRelationPlatformView)))))) // ?anime=<anime_id>&platform=<platform_id>
	router.Handle("/admin/relation/plat/update/", base.BaseMiddleware(base.PutOnlyMiddleware(base.AdminRequiredMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.UpdateRelationPlatformView))))))

	// company
	router.Handle("/admin/company/post/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.PostOnlyMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.InsertCompanyView))))))
	router.Handle("/admin/company/edit/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.PutOnlyMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.UpdateCompanyView))))))
	router.Handle("/admin/company/", base.BaseMiddleware(base.AdminRequiredMiddlewareGet(http.HandlerFunc(adminController.AdminCompanyView))))

	// staff
	router.Handle("/admin/staff/post/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.PostOnlyMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.InsertStaffView))))))

	// role
	router.Handle("/admin/role/", base.BaseMiddleware(base.AdminRequiredMiddleware(http.HandlerFunc(adminController.ListRoleView))))
	router.Handle("/admin/role/post/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.PostOnlyMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.InsertRoleView))))))
	router.Handle("/admin/staffrole/post/", base.BaseMiddleware(base.AdminRequiredMiddleware(base.PostOnlyMiddleware(base.VerifyCsrfMiddleware(http.HandlerFunc(adminController.InsertStaffRoleView))))))

	// test for csrf
	router.Handle("/csrf/", base.BaseMiddleware(http.HandlerFunc(base.SetCsrfCookieView)))
	// router.Handle("/test/verify/", base.BaseMiddleware(base.VerifyCsrfMiddleware(
	// 	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprint(w, cache.Items[domain.CacheTypeCsrf])
	// 	}),
	// )))

	if tools.IsProductionEnv() {
		if err := http.ListenAndServe(":8000", infrastructure.Log(router, lg)); err != nil {
			domain.ErrorAlert(err)
		}
	} else {
		if err := http.ListenAndServe(":8000", infrastructure.Log(router, lg)); err != nil {
			domain.ErrorAlert(err)
		}
	}

}
