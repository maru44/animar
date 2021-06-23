package usecase

import "animar/v1/pkg/domain"

type BlogInteractor struct {
	repository BlogRepository
}

func NewBlogInteractor(blog BlogRepository) domain.BlogInteractor {
	return &BlogInteractor{
		repository: blog,
	}
}

/************************
        repository
************************/

type BlogRepository interface {
	ListAll() (domain.TBlogs, error)
	FilterByUser(string, string) (domain.TBlogs, error)
	GetUserId(int) (string, error)
	FindById(int) (domain.TBlog, error)
	FindBySlug(string) (domain.TBlog, error)
	Insert(domain.TBlogInsert) (int, error)
	Update(domain.TBlogInsert, int) (int, error)
	Delete(int) (int, error)
	// relation
	FilterByBlog(int) ([]domain.TJoinedAnime, error)
	// FilterByAnime(int) (domain.TBlogJoinAnimes, error)
	InsertRelation(int, int) (bool, error)
	DeleteRelation(int, int) error
}

/**********************
   interactor methods
***********************/

func (interactor *BlogInteractor) ListBlog() (blogs domain.TBlogs, err error) {
	blogs, err = interactor.repository.ListAll()
	return
}

func (interactor *BlogInteractor) ListBlogByUser(accessUserId string, blogUserId string) (blogs domain.TBlogs, err error) {
	blogs, err = interactor.repository.FilterByUser(accessUserId, blogUserId)
	return
}

func (interactor *BlogInteractor) BlogUserId(blogId int) (userId string, err error) {
	userId, err = interactor.repository.GetUserId(blogId)
	return
}

func (interactor *BlogInteractor) DetailBlog(id int) (blog domain.TBlog, err error) {
	blog, err = interactor.repository.FindById(id)
	return
}

func (interactor *BlogInteractor) DetailBlogBySlug(slug string) (blog domain.TBlog, err error) {
	blog, err = interactor.repository.FindBySlug(slug)
	return
}

func (interactor *BlogInteractor) InsertBlog(blog domain.TBlogInsert) (lastInserted int, err error) {
	lastInserted, err = interactor.repository.Insert(blog)
	return
}

func (interactor *BlogInteractor) UpdateBlog(blog domain.TBlogInsert, id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.repository.Update(blog, id)
	return
}

func (interactor *BlogInteractor) DeleteBlog(id int) (rowsAffected int, err error) {
	rowsAffected, err = interactor.repository.Delete(id)
	return
}

func (interactor *BlogInteractor) RelationAnimeByBlog(blogId int) (animes []domain.TJoinedAnime, err error) {
	animes, err = interactor.repository.FilterByBlog(blogId)
	return
}

func (interactor *BlogInteractor) InsertRelationAnime(animeId int, blogId int) (is_success bool, err error) {
	is_success, err = interactor.repository.InsertRelation(animeId, blogId)
	return
}

func (interactor *BlogInteractor) DeleteRelationAnime(animeId int, blogId int) (err error) {
	err = interactor.repository.DeleteRelation(animeId, blogId)
	return
}
