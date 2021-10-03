package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"net/http"

	"github.com/maru44/perr"
)

type StaffController struct {
	interactor domain.StaffInteractor
}

func NewStaffController(sqlHandler database.SqlHandler) *StaffController {
	return &StaffController{
		interactor: usecase.NewStaffInteractor(
			&database.StaffRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (sfc *StaffController) StaffListView(w http.ResponseWriter, r *http.Request) {
	staffs, err := sfc.interactor.StaffList()
	response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": staffs})
	return
}
