package usecase

import "animar/v1/internal/pkg/domain"

type CompanyInteractor struct {
	repo CompanyRepository
}

/************************
        repository
************************/

type CompanyRepository interface {
	List() ([]domain.Company, error)
	DetailByEng(string) (domain.CompanyDetail, error)
	// admin
	Insert(domain.CompanyInput) (int, error)
	Update(domain.CompanyInput, string) (int, error)
}

/**********************
   interactor methods
***********************/
