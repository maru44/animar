package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/mvc/auth"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/usecase"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type BlogController struct {
	interactor domain.BlogInteractor
	BaseController
}

func NewBlogController(sqlHandler database.SqlHandler) *BlogController {
	return &BlogController{
		interactor: usecase.NewBlogInteractor(
			&database.BlogRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (controller *BlogController) BlogListView(w http.ResponseWriter, r *http.Request) (ret error) {
	blogs, err := controller.interactor.ListBlog()
	ret = response(w, err, map[string]interface{}{"data": blogs})
	return ret
}

func (controller *BlogController) BlogJoinAnimeView(w http.ResponseWriter, r *http.Request) (ret error) {
	var userId string
	switch r.Method {
	case "GET":
		userId, _ = controller.getUserIdFromCookie(r)
	case "POST":
		var posted fire.TUserIdCookieInput
		json.NewDecoder(r.Body).Decode(&posted)
		userId = fire.GetUserIdFromToken(posted.Token)
	default:
		userId = ""
	}

	query := r.URL.Query()
	slug := query.Get("s")
	id := query.Get("id")
	uid := query.Get("u")

	if slug != "" {
		blog, err := controller.interactor.DetailBlogBySlug(slug)
		blog.Animes, _ = controller.interactor.RelationAnimeByBlog(blog.GetId())
		ret = response(w, err, map[string]interface{}{"data": blog})
	} else if id != "" {
		i, _ := strconv.Atoi(id)
		blog, err := controller.interactor.DetailBlog(i)
		blog.Animes, _ = controller.interactor.RelationAnimeByBlog(i)
		ret = response(w, err, map[string]interface{}{"data": blog})
	} else if uid != "" {
		blogs, err := controller.interactor.ListBlogByUser(userId, uid)
		ret = response(w, err, map[string]interface{}{"data": blogs})
	} else {
		blogs, err := controller.interactor.ListBlog()
		ret = response(w, err, map[string]interface{}{"data": blogs})
	}
	return ret
}

func (controller *BlogController) InsertBlogWithRelationView(w http.ResponseWriter, r *http.Request) (ret error) {
	/*  userId 取得  */
	idToken, _ := r.Cookie("idToken")
	claims := auth.VerifyFirebase(context.Background(), idToken.Value)
	userId := claims["user_id"].(string)

	if userId == "" {
		ret = response(w, domain.ErrUnauthorized, nil)
	} else {
		var p domain.TBlogInsert
		json.NewDecoder(r.Body).Decode(&p)
		lastInserted, err := controller.interactor.InsertBlog(p, userId)
		for _, animeId := range p.AnimeIds {
			controller.interactor.InsertRelationAnime(animeId, lastInserted)
		}
		ret = response(w, err, map[string]interface{}{"data": lastInserted})
	}
	return ret
}

func (controller *BlogController) UpdateBlogWithRelationView(w http.ResponseWriter, r *http.Request) (ret error) {
	/*  userId 取得  */
	idToken, _ := r.Cookie("idToken")
	claims := auth.VerifyFirebase(context.Background(), idToken.Value)
	userId := claims["user_id"].(string)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId, _ := controller.interactor.BlogUserId(id)
	if blogUserId != userId {
		ret = response(w, domain.ErrForbidden, nil)
	} else {
		var p domain.TBlogInsert
		json.NewDecoder(r.Body).Decode(&p)
		rowsAffected, err := controller.interactor.UpdateBlog(p, id)
		ret = response(w, err, map[string]interface{}{"data": rowsAffected})
	}
	return ret
}

func (controller *BlogController) DeleteBlogView(w http.ResponseWriter, r *http.Request) (ret error) {
	/*  userId 取得  */
	idToken, _ := r.Cookie("idToken")
	claims := auth.VerifyFirebase(context.Background(), idToken.Value)
	userId := claims["user_id"].(string)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId, _ := controller.interactor.BlogUserId(id)
	if blogUserId != userId {
		ret = response(w, domain.ErrForbidden, nil)
	} else {
		deletedRow, err := controller.interactor.DeleteBlog(id)
		ret = response(w, err, map[string]interface{}{"data": deletedRow})
	}
	return ret
}
