package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/usecase"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type BlogController struct {
	interactor domain.BlogInteractor
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

func (controller *BlogController) BlogListView(w http.ResponseWriter, r *http.Request) error {
	blogs, _ := controller.interactor.ListBlog()
	api.JsonResponse(w, map[string]interface{}{"data": blogs})
	return nil
}

func (controller *BlogController) BlogJoinAnimeView(w http.ResponseWriter, r *http.Request) error {
	var userId string
	switch r.Method {
	case "GET":
		userId = fire.GetIdFromCookie(r)
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
		blog, _ := controller.interactor.DetailBlogBySlug(slug)
		blog.Animes, _ = controller.interactor.RelationAnimeByBlog(blog.GetId())
		api.JsonResponse(w, map[string]interface{}{"data": blog})
	} else if id != "" {
		i, _ := strconv.Atoi(id)
		blog, _ := controller.interactor.DetailBlog(i)
		blog.Animes, _ = controller.interactor.RelationAnimeByBlog(i)
		api.JsonResponse(w, map[string]interface{}{"data": blog})
	} else if uid != "" {
		blogs, _ := controller.interactor.ListBlogByUser(userId, uid)
		api.JsonResponse(w, map[string]interface{}{"data": blogs})
	} else {
		blogs, _ := controller.interactor.ListBlog()
		api.JsonResponse(w, map[string]interface{}{"data": blogs})
	}
	return nil
}

func (controller *BlogController) InsertBlogWithRelationView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized")
	} else {
		var p domain.TBlogInsert
		json.NewDecoder(r.Body).Decode(&p)
		lastInserted, _ := controller.interactor.InsertBlog(p, userId)
		for _, animeId := range p.AnimeIds {
			controller.interactor.InsertRelationAnime(animeId, lastInserted)
		}
		api.JsonResponse(w, map[string]interface{}{"data": lastInserted})
	}
	return nil
}

func (controller *BlogController) UpdateBlogWithRelationView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId, _ := controller.interactor.BlogUserId(id)
	if blogUserId != userId {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("Forbidden")
	}
	var p domain.TBlogInsert
	json.NewDecoder(r.Body).Decode(&p)
	rowsAffected, _ := controller.interactor.UpdateBlog(p, id)
	api.JsonResponse(w, map[string]interface{}{"data": rowsAffected})
	return nil
}

func (controller *BlogController) DeleteBlogView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId, _ := controller.interactor.BlogUserId(id)
	if blogUserId != userId {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("Forbidden")
	}
	deletedRow, _ := controller.interactor.DeleteBlog(id)
	api.JsonResponse(w, map[string]interface{}{"data": deletedRow})
	return nil
}
