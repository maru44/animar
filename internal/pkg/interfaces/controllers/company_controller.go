package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"net/http"
)

type CompanyController struct {
	coi domain.CompanyInteractor
}

func NewCompanyController(sqlHandler database.SqlHandler) *CompanyController {
	return &CompanyController{
		coi: usecase.NewCompanyInteractor(
			&database.CompanyRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (coc *CompanyController) ListCompanyView(w http.ResponseWriter, r *http.Request) {
	cs, err := coc.coi.ListCompany()
	response(w, err, map[string]interface{}{"data": cs})
	return
}
