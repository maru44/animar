package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/usecase"
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
		BaseController: *NewBaseController(),
	}
}

func (controller *BlogController) BlogListView(w http.ResponseWriter, r *http.Request) {
	blogs, err := controller.interactor.ListBlog()
	response(w, err, map[string]interface{}{"data": blogs})
	return
}

func (controller *BlogController) BlogJoinAnimeView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	query := r.URL.Query()
	slug := query.Get("s")
	id := query.Get("id")
	uid := query.Get("u")

	if slug != "" {
		blog, err := controller.interactor.DetailBlogBySlug(slug)
		blog.Animes, _ = controller.interactor.RelationAnimeByBlog(blog.GetId())
		response(w, err, map[string]interface{}{"data": blog})
	} else if id != "" {
		i, _ := strconv.Atoi(id)
		blog, err := controller.interactor.DetailBlog(i)
		blog.Animes, _ = controller.interactor.RelationAnimeByBlog(i)
		response(w, err, map[string]interface{}{"data": blog})
	} else if uid != "" {
		blogs, err := controller.interactor.ListBlogByUser(userId, uid)
		response(w, err, map[string]interface{}{"data": blogs})
	} else {
		blogs, err := controller.interactor.ListBlog()
		response(w, err, map[string]interface{}{"data": blogs})
	}
	return
}

func (controller *BlogController) InsertBlogWithRelationView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)

	var p domain.TBlogInsert
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := controller.interactor.InsertBlog(p, userId)
	for _, animeId := range p.AnimeIds {
		controller.interactor.InsertRelationAnime(animeId, lastInserted)
	}
	response(w, err, map[string]interface{}{"data": lastInserted})
	return
}

func (controller *BlogController) UpdateBlogWithRelationView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId, _ := controller.interactor.BlogUserId(id)
	if blogUserId != userId {
		response(w, domain.ErrForbidden, nil)
	} else {
		var p domain.TBlogInsert
		json.NewDecoder(r.Body).Decode(&p)
		rowsAffected, err := controller.interactor.UpdateBlog(p, id)
		response(w, err, map[string]interface{}{"data": rowsAffected})
	}
	return
}

func (controller *BlogController) DeleteBlogView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId, _ := controller.interactor.BlogUserId(id)
	if blogUserId != userId {
		response(w, domain.ErrForbidden, nil)
	} else {
		deletedRow, err := controller.interactor.DeleteBlog(id)
		response(w, err, map[string]interface{}{"data": deletedRow})
	}
	return
}
