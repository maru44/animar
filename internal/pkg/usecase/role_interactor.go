package usecase

import "animar/v1/internal/pkg/domain"

type RoleInteractor struct {
	repo RoleRepository
}

func NewRoleInteractor(sfr RoleRepository) domain.RoleInteractor {
	return &RoleInteractor{
		repo: sfr,
	}
}

/************************
        repository
************************/

type RoleRepository interface {
	FilterByAnime(int) ([]domain.AnimeStaffRole, error)
	/*  admin  */
	List() ([]domain.Role, error)
	Insert(domain.RoleInput) (int, error)
	InsertStaffRole(domain.AnimeStaffRoleInput) (int, error)
}

/************************
        methods
************************/

func (ri *RoleInteractor) StaffRoleByAnime(animeId int) ([]domain.AnimeStaffRole, error) {
	return ri.repo.FilterByAnime(animeId)
}
