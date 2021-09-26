package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/interfaces/s3"
	"animar/v1/internal/pkg/usecase"
	"encoding/json"
	"net/http"
	"strconv"
)

type BlogController struct {
	interactor domain.BlogInteractor
	s3         domain.S3Interactor
}

func NewBlogController(sqlHandler database.SqlHandler, uploader s3.Uploader) *BlogController {
	return &BlogController{
		interactor: usecase.NewBlogInteractor(
			&database.BlogRepository{
				SqlHandler: sqlHandler,
			},
		),
		s3: usecase.NewS3Interactor(
			&s3.S3Repository{
				Uploader: uploader,
			},
		),
	}
}

func (controller *BlogController) BlogListView(w http.ResponseWriter, r *http.Request) {
	blogs, err := controller.interactor.ListBlog()
	response(w, r, err, map[string]interface{}{"data": blogs})
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
		response(w, r, err, map[string]interface{}{"data": blog})
	} else if id != "" {
		i, _ := strconv.Atoi(id)
		blog, err := controller.interactor.DetailBlog(i)
		blog.Animes, _ = controller.interactor.RelationAnimeByBlog(i)
		response(w, r, err, map[string]interface{}{"data": blog})
	} else if uid != "" {
		blogs, err := controller.interactor.ListBlogByUser(userId, uid)
		response(w, r, err, map[string]interface{}{"data": blogs})
	} else {
		blogs, err := controller.interactor.ListBlog()
		response(w, r, err, map[string]interface{}{"data": blogs})
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
	response(w, r, err, map[string]interface{}{"data": lastInserted})
	return
}

func (controller *BlogController) UpdateBlogWithRelationView(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(USER_ID).(string)
	query := r.URL.Query()
	strId := query.Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, domain.NewWrapError(err, domain.ExternalServerError), nil)
		return
	}

	// user 不一致
	blogUserId, _ := controller.interactor.BlogUserId(id)
	if blogUserId != userId {
		response(w, r, domain.NewError("forbidden", domain.ForbiddenError), nil)
	} else {
		var p domain.TBlogInsert
		json.NewDecoder(r.Body).Decode(&p)
		rowsAffected, err := controller.interactor.UpdateBlog(p, id)
		response(w, r, err, map[string]interface{}{"data": rowsAffected})
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
		response(w, r, domain.NewError("forbidden", domain.ForbiddenError), nil)
	} else {
		deletedRow, err := controller.interactor.DeleteBlog(id)
		response(w, r, err, map[string]interface{}{"data": deletedRow})
	}
	return
}

/*************************
        image upload
************************/

func (controller *BlogController) SimpleUploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB
	var returnFileName string

	file, fileHeader, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		returnFileName, err = controller.s3.Image(file, fileHeader.Filename, []string{"column", "content"})
	}
	response(w, r, err, map[string]interface{}{"data": returnFileName})
	return
}
