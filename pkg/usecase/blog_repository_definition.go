package usecase

import "animar/v1/pkg/domain"

type BlogRepository interface {
	ListAll() (domain.TBlogs, error)
	FilterByUser(string, string) (domain.TBlogs, error)
	GetUserId(int) (string, error)
	FilterByAnime()
}
