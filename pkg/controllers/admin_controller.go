package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/s3"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"encoding/json"
	"net/http"
	"strconv"
)

type AdminController struct {
	AnimeInteractor usecase.AdminAnimeInteractor
}

func NewAdminController(sqlHandler database.SqlHandler) *AdminController {
	return &AdminController{
		AnimeInteractor: usecase.AdminAnimeInteractor{
			AdminAnimeRepository: &database.AdminRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *AdminController) AnimeListAdminView(w http.ResponseWriter, r *http.Request) error {
	animes, _ := controller.AnimeInteractor.AnimesAllAdmin()
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}

func (controller *AdminController) AnimeDetailAdminView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	var posted api.TUserToken
	json.NewDecoder(r.Body).Decode(&posted)

	anime, err := controller.AnimeInteractor.AnimeDetailAdmin(id)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": anime})
	return nil
}

func (controller *AdminController) AnimePostAdminView(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB
	var returnFileName string
	var insertedId int
	file, fileHeader, err := r.FormFile("thumb")
	if err == nil { // w/ thumb picture
		defer file.Close()
		returnFileName, err = s3.UploadS3(file, fileHeader.Filename, []string{"anime"})

		if err != nil {
			tools.ErrorLog(err)
			return err
		}
	} else { // w/o thumb picture
		returnFileName = r.FormValue("pre_thumb")
	}
	slug := tools.GenRandSlug(12)
	series, _ := strconv.Atoi(r.FormValue("series_id"))
	episodes, _ := strconv.Atoi(r.FormValue("count_episodes"))
	a := domain.TAnimeInsert{
		Title:         r.FormValue("title"),
		Slug:          slug,
		Kana:          tools.NewNullString("kana"),
		EngName:       tools.NewNullString(r.FormValue("eng_name")),
		Abbreviation:  tools.NewNullString(r.FormValue("abbreviation")),
		Description:   tools.NewNullString(r.FormValue("description")),
		State:         tools.NewNullString(r.FormValue("state")),
		ThumbUrl:      tools.NewNullString(returnFileName),
		SeriesId:      tools.NewNullInt(series),
		CountEpisodes: tools.NewNullInt(episodes),
		Copyright:     tools.NewNullString(r.FormValue("copyright")),
	}
	insertedId, err = controller.AnimeInteractor.AnimeInsert(a)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": insertedId})
	return nil
}

// @TODO edit & delete
