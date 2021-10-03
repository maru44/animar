package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"net/http"
	"strconv"

	"github.com/maru44/perr"
)

type RoleController struct {
	interactor domain.RoleInteractor
}

func NewRoleController(sqlHandler database.SqlHandler) *RoleController {
	return &RoleController{
		interactor: usecase.NewRoleInteractor(
			&database.RoleRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (controller *RoleController) ListStaffRoleView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("anime")
	animeId, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	staffRoles, err := controller.interactor.StaffRoleByAnime(animeId)
	response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": staffRoles})
	return
}
