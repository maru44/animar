package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/interfaces/database"
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/fire"
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
				SqlHandler: sqlHandler,
			},
			&database.AdminSeriesRepository{
				SqlHandler: sqlHandler,
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

func (controller *AdminController) PlatformUpdateView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

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
	rowsAffected, err := controller.interactor.PlatformUpdate(p, id)
	api.JsonResponse(w, map[string]interface{}{"data": rowsAffected})
	return nil
}

func (controller *AdminController) PlatformDeleteview(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	rowsAffected, _ := controller.interactor.PlatformDelete(id)
	api.JsonResponse(w, map[string]interface{}{"data": rowsAffected})
	return nil
}

func (controller *AdminController) RelationPlatformView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id") // animeId
	id, _ := strconv.Atoi(strId)
	relations, _ := controller.interactor.RelationPlatformByAnime(id)
	api.JsonResponse(w, map[string]interface{}{"data": relations})
	return nil
}

func (controller *AdminController) InsertRelationPlatformView(w http.ResponseWriter, r *http.Request) error {
	var p domain.TRelationPlatformInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := controller.interactor.RelationPlatformInsert(p)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": lastInserted})
	return nil
}

func (controller *AdminController) DeleteRelationPlatformView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strAnimeId := query.Get("anime")
	strPlatformId := query.Get("platform")
	animeId, _ := strconv.Atoi(strAnimeId)
	platformId, _ := strconv.Atoi(strPlatformId)

	rowsAffected, err := controller.interactor.RelationPlatformDelete(animeId, platformId)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": rowsAffected})
	return nil
}

/************************
         season
*************************/

func (controller *AdminController) SeasonView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, _ := strconv.Atoi(id)
		s, _ := controller.interactor.DetailSeason(i)
		api.JsonResponse(w, map[string]interface{}{"data": s})
	} else {
		s, _ := controller.interactor.ListSeason()
		api.JsonResponse(w, map[string]interface{}{"data": s})
	}
	return nil
}

func (controller *AdminController) InsertSeasonView(w http.ResponseWriter, r *http.Request) error {
	var p domain.TSeasonInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := controller.interactor.InsertSeason(p)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": lastInserted})
	return nil
}

func (controller *AdminController) InsertRelationSeasonView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetAdminIdFromCookie(r)

	if userId == "" {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("Forbidden")
	} else {
		var s domain.TSeasonRelationInput
		json.NewDecoder(r.Body).Decode(&s)
		lastInserted, err := controller.interactor.InsertRelationSeasonAnime(s)
		if err != nil {
			tools.ErrorLog(err)
			return err
		}
		api.JsonResponse(w, map[string]interface{}{"data": lastInserted})
	}
	return nil
}

/************************
         series
*************************/

func (controller *AdminController) SeriesView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, _ := strconv.Atoi(id)
		s, _ := controller.interactor.DetailSeries(i)
		api.JsonResponse(w, map[string]interface{}{"data": s})
	} else {
		s, _ := controller.interactor.ListSeries()
		api.JsonResponse(w, map[string]interface{}{"data": s})
	}
	return nil
}

func (controller *AdminController) InsertSeriesView(w http.ResponseWriter, r *http.Request) error {
	var p domain.TSeriesInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := controller.interactor.InsertSeries(p)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": lastInserted})
	return nil
}

func (controller *AdminController) UpdateSeriesView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	var p domain.TSeriesInput
	json.NewDecoder(r.Body).Decode(&p)
	rowsAffected, err := controller.interactor.UpdateSeries(p, id)
	if err != nil {
		tools.ErrorLog(err)
		return err
	}
	api.JsonResponse(w, map[string]interface{}{"data": rowsAffected})
	return nil
}
