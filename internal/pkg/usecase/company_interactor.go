package usecase

import "animar/v1/internal/pkg/domain"

type CompanyInteractor struct {
	repo CompanyRepository
}

func NewCompanyInteractor(cor CompanyRepository) domain.CompanyInteractor {
	return &CompanyInteractor{
		repo: cor,
	}
}

/************************
        repository
************************/

type CompanyRepository interface {
	Insert(domain.CompanyInput) (int, error)
	List() ([]domain.Company, error)
}

/**********************
   interactor methods
***********************/

func (ci *CompanyInteractor) ListCompany() ([]domain.Company, error) {
	return ci.repo.List()
}
