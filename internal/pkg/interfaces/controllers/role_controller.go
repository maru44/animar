package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"net/http"
	"strconv"
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
	animeId, _ := strconv.Atoi(strId)

	staffRoles, err := controller.interactor.StaffRoleByAnime(animeId)
	response(w, r, err, map[string]interface{}{"data": staffRoles})
	return
}
