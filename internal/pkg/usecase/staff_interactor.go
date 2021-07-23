package usecase

import "animar/v1/internal/pkg/domain"

type StaffInteractor struct {
	repo StaffRepository
}

func NewStaffInteractor(sfr StaffRepository) domain.StaffInteractor {
	return &StaffInteractor{
		repo: sfr,
	}
}

/************************
        repository
************************/

type StaffRepository interface {
	List() ([]domain.Staff, error)
	Insert(domain.StaffInput) (int, error)
}

/************************
        methods
************************/

func (sfi *StaffInteractor) StaffList() ([]domain.Staff, error) {
	return sfi.repo.List()
}
