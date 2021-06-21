package usecase

import "animar/v1/pkg/domain"

type AdminAnimeRepository interface {
	ListAll() (domain.TAnimes, error)
	FindById(int) (domain.TAnimeAdmin, error)
	Insert(domain.TAnimeInsert) (int, error)
	Update(int, domain.TAnimeInsert) (int, error)
	Delete(int) (int, error)
}
