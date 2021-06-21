package usecase

import "animar/v1/pkg/domain"

type AdminAnimeRepository interface {
	ListAll() (domain.TAnimes, error)
	FindById(int) (domain.TAnimeAdmin, error)
	Insert(domain.TAnimeInsert) (int, error)
	Update(int, domain.TAnimeInsert) (int, error)
	Delete(int) (int, error)
}

type AdminPlatformRepository interface {
	ListAll() (domain.TPlatforms, error)
	Insert(domain.TPlatform) (int, error)
	FindById(int) (domain.TPlatform, error)
	Update(domain.TPlatform, int) (int, error)
	Delete(int) (int, error)
	// relation
	FilterByAnime(int) (domain.TRelationPlatforms, error)
	InsertRelation(domain.TRelationPlatformInput) (int, error)
	DeleteRelation(int, int) (int, error)
}
