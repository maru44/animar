package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/usecase"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/maru44/perr"
)

type PlatformController struct {
	interactor domain.PlatformInteractor
}

func NewPlatformController(sqlHandler database.SqlHandler) *PlatformController {
	return &PlatformController{
		interactor: usecase.NewPlatformInteractor(
			&database.PlatformRepository{
				SqlHandler: sqlHandler,
			},
		),
	}
}

func (con *PlatformController) RelationPlatformByAnimeView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
	}

	platforms, err := con.interactor.RelationPlatformByAnime(id)
	response(w, r, perr.Wrap(err, perr.NotFound), map[string]interface{}{"data": platforms})
	return
}

func (con *PlatformController) RegisterNotifiedTargetView(w http.ResponseWriter, r *http.Request) {
	var p domain.NotifiedTargetInput
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	inserted, err := con.interactor.RegisterNotifiedTarget(p)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	response(w, r, nil, map[string]interface{}{"data": inserted})
	return
}
