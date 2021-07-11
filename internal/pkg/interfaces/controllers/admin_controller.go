package controllers

import (
	"animar/v1/internal/pkg/domain"
	"animar/v1/internal/pkg/interfaces/database"
	"animar/v1/internal/pkg/interfaces/s3"
	"animar/v1/internal/pkg/tools/tools"
	"animar/v1/internal/pkg/usecase"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type AdminController struct {
	interactor domain.AdminInteractor
	s3         domain.S3Interactor
}

func NewAdminController(sqlHandler database.SqlHandler, uploader s3.Uploader) *AdminController {
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
		s3: usecase.NewS3Interactor(
			&s3.S3Repository{
				Uploader: uploader,
			},
		),
	}
}

/************************
         anime
*************************/

func (controller *AdminController) AnimeListAdminView(w http.ResponseWriter, r *http.Request) {
	animes, err := controller.interactor.AnimesAllAdmin()
	response(w, err, map[string]interface{}{"data": animes})
	return
}

func (controller *AdminController) AnimeDetailAdminView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	var posted domain.TUserToken
	json.NewDecoder(r.Body).Decode(&posted)

	anime, err := controller.interactor.AnimeDetailAdmin(id)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": anime})
	return
}

func (controller *AdminController) AnimePostAdminView(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB
	var returnFileName string
	var insertedId int
	file, fileHeader, err := r.FormFile("thumb")
	if err == nil { // w/ thumb picture
		defer file.Close()
		returnFileName, err = controller.s3.Image(file, fileHeader.Filename, []string{"anime"})

		if err != nil {
			tools.ErrorLog(err)
			response(w, err, nil)
			return
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
		Kana:          tools.NewNullString(r.FormValue("kana")),
		EngName:       tools.NewNullString(r.FormValue("engName")),
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
	}
	response(w, err, map[string]interface{}{"data": insertedId})
	return
}

func (controller *AdminController) AnimeUpdateView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("thumb")
	var returnFileName string
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = controller.s3.Image(file, fileHeader.Filename, []string{"anime"})
		if err != nil {
			tools.ErrorLog(err)
			response(w, err, nil)
			return
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
	}
	response(w, err, map[string]interface{}{"data": rowsAffected})
	return
}

func (controller *AdminController) AnimeDeleteView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	rowsAffected, err := controller.interactor.AnimeDelete(id)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": rowsAffected})
	return
}

/************************
         platform
*************************/

func (controller *AdminController) PlatformView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, _ := strconv.Atoi(id)
		platform, err := controller.interactor.PlatformDetail(i)
		if err != nil || platform.ID == 0 {
			err = domain.ErrNotFound
		}
		response(w, err, map[string]interface{}{"data": platform})
	} else {
		platforms, err := controller.interactor.PlatformAllAdmin()
		response(w, err, map[string]interface{}{"data": platforms})
	}
	return
}

func (controller *AdminController) PlatformInsertView(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("image")
	var returnFileName string
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = controller.s3.Image(file, fileHeader.Filename, []string{"platform"})

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
	response(w, err, map[string]interface{}{"data": lastInserted})
	return
}

func (controller *AdminController) PlatformUpdateView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("image")
	var returnFileName string
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = controller.s3.Image(file, fileHeader.Filename, []string{"platform"})

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
	response(w, err, map[string]interface{}{"data": rowsAffected})
	return
}

func (controller *AdminController) PlatformDeleteview(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	rowsAffected, err := controller.interactor.PlatformDelete(id)
	response(w, err, map[string]interface{}{"data": rowsAffected})
	return
}

func (controller *AdminController) RelationPlatformView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id") // animeId
	id, _ := strconv.Atoi(strId)
	relations, err := controller.interactor.RelationPlatformByAnime(id)
	response(w, err, map[string]interface{}{"data": relations})
	return
}

func (controller *AdminController) InsertRelationPlatformView(w http.ResponseWriter, r *http.Request) {
	var p domain.TRelationPlatformInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := controller.interactor.RelationPlatformInsert(p)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": lastInserted})
	return
}

func (controller *AdminController) DeleteRelationPlatformView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strAnimeId := query.Get("anime")
	strPlatformId := query.Get("platform")
	animeId, _ := strconv.Atoi(strAnimeId)
	platformId, _ := strconv.Atoi(strPlatformId)

	rowsAffected, err := controller.interactor.RelationPlatformDelete(animeId, platformId)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": rowsAffected})
	return
}

/************************
         season
*************************/

func (controller *AdminController) SeasonView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, _ := strconv.Atoi(id)
		s, err := controller.interactor.DetailSeason(i)
		response(w, err, map[string]interface{}{"data": s})
	} else {
		s, err := controller.interactor.ListSeason()
		response(w, err, map[string]interface{}{"data": s})
	}
	return
}

func (controller *AdminController) InsertSeasonView(w http.ResponseWriter, r *http.Request) {
	var p domain.TSeasonInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := controller.interactor.InsertSeason(p)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": lastInserted})
	return
}

func (controller *AdminController) InsertRelationSeasonView(w http.ResponseWriter, r *http.Request) {
	var s domain.TSeasonRelationInput
	json.NewDecoder(r.Body).Decode(&s)
	lastInserted, err := controller.interactor.InsertRelationSeasonAnime(s)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": lastInserted})
	return
}

/************************
         series
*************************/

func (controller *AdminController) SeriesView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, _ := strconv.Atoi(id)
		s, err := controller.interactor.DetailSeries(i)
		response(w, err, map[string]interface{}{"data": s})
	} else {
		s, err := controller.interactor.ListSeries()
		response(w, err, map[string]interface{}{"data": s})
	}
	return
}

func (controller *AdminController) InsertSeriesView(w http.ResponseWriter, r *http.Request) {
	var p domain.TSeriesInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := controller.interactor.InsertSeries(p)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": lastInserted})
	return
}

func (controller *AdminController) UpdateSeriesView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	var p domain.TSeriesInput
	json.NewDecoder(r.Body).Decode(&p)
	rowsAffected, err := controller.interactor.UpdateSeries(p, id)
	if err != nil {
		tools.ErrorLog(err)
	}
	response(w, err, map[string]interface{}{"data": rowsAffected})
	return
}
