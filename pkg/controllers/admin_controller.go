package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/s3"
	"animar/v1/pkg/tools/tools"
	"animar/v1/pkg/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type AdminController struct {
	interactor domain.AdminInteractor
}

func NewAdminController(sqlHandler database.SqlHandler) *AdminController {
	return &AdminController{
		interactor: usecase.NewAdminAnimeInteractor(
			&database.AdminAnimeRepository{
				SqlHandler: sqlHandler,
			},
			&database.AdminPlatformRepository{
				SqlHandler: sqlHandler,
			},
			&database.AdminSeasonRepository{
				sqlHandler: sqlHandler,
			},
		),
	}
}

/************************
         anime
*************************/

func (controller *AdminController) AnimeListAdminView(w http.ResponseWriter, r *http.Request) error {
	animes, _ := controller.interactor.AnimesAllAdmin()
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}

func (controller *AdminController) AnimeDetailAdminView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	var posted api.TUserToken
	json.NewDecoder(r.Body).Decode(&posted)

	anime, err := controller.interactor.AnimeDetailAdmin(id)
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
	insertedId, err = controller.interactor.AnimeInsert(a)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": insertedId})
	return nil
}

func (controller *AdminController) AnimeUpdateView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("thumb")
	var returnFileName string
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = s3.UploadS3(file, fileHeader.Filename, []string{"anime"})
		if err != nil {
			tools.ErrorLog(err)
		}
	} else {
		returnFileName = r.FormValue("pre_thumb")
	}
	series, _ := strconv.Atoi(r.FormValue("series_id"))
	episodes, _ := strconv.Atoi(r.FormValue("count_episodes"))
	a := domain.TAnimeInsert{
		Title:         r.FormValue("title"),
		Abbreviation:  tools.NewNullString(r.FormValue("abbreviation")),
		Kana:          tools.NewNullString(r.FormValue("kana")),
		EngName:       tools.NewNullString(r.FormValue("engName")),
		Description:   tools.NewNullString(r.FormValue("description")),
		State:         tools.NewNullString(r.FormValue("state")),
		SeriesId:      tools.NewNullInt(series),
		CountEpisodes: tools.NewNullInt(episodes),
		Copyright:     tools.NewNullString(r.FormValue("copyright")),
		ThumbUrl:      tools.NewNullString(returnFileName),
	}

	rowsAffected, err := controller.interactor.AnimeUpdate(id, a)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": rowsAffected})
	return nil
}

func (controller *AdminController) AnimeDeleteView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	rowsAffected, err := controller.interactor.AnimeDelete(id)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": rowsAffected})
	return nil
}

/************************
         platform
*************************/

// @TODO platform endpoint

func (controller *AdminController) PlatformView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, _ := strconv.Atoi(id)
		platform, err := controller.interactor.PlatformDetail(i)
		if err != nil || platform.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": platform})
	} else {
		platforms, _ := controller.interactor.PlatformAllAdmin()
		api.JsonResponse(w, map[string]interface{}{"data": platforms})
	}
	return nil
}

func (controller *AdminController) PlatformInsertView(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("image")
	var returnFileName string
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = s3.UploadS3(file, fileHeader.Filename, []string{"platform"})

		if err != nil {
			fmt.Print(err)
		}
	} else {
		returnFileName = ""
	}
	validStr := r.FormValue("valid")
	isValid, _ := strconv.ParseBool(validStr)

	p := domain.TPlatform{
		EngName:  r.FormValue("engName"),
		PlatName: tools.NewNullString(r.FormValue("platName")),
		BaseUrl:  tools.NewNullString(r.FormValue("baseUrl")),
		Image:    tools.NewNullString(returnFileName),
		IsValid:  isValid,
	}
	lastInserted, err := controller.interactor.PlatformInsert(p)
	api.JsonResponse(w, map[string]interface{}{"data": lastInserted})
	return nil
}
