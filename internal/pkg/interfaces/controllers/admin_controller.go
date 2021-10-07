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

	"github.com/maru44/perr"
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
			&database.CompanyRepository{
				SqlHandler: sqlHandler,
			},
			&database.StaffRepository{
				SqlHandler: sqlHandler,
			},
			&database.RoleRepository{
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

func (adc *AdminController) AnimeListAdminView(w http.ResponseWriter, r *http.Request) {
	animes, err := adc.interactor.AnimesAllAdmin()
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": animes})
	return
}

func (adc *AdminController) AnimeDetailAdminView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	var posted domain.TUserToken
	json.NewDecoder(r.Body).Decode(&posted)

	anime, err := adc.interactor.AnimeDetailAdmin(id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": anime})
	return
}

func (adc *AdminController) AnimePostAdminView(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB
	var returnFileName string
	var insertedId int
	file, fileHeader, err := r.FormFile("thumb")
	if err == nil { // w/ thumb picture
		defer file.Close()
		returnFileName, err = adc.s3.Image(file, fileHeader.Filename, []string{"anime"})

		if err != nil {
			response(w, r, perr.Wrap(err, perr.BadRequest), nil)
			return
		}
	} else { // w/o thumb picture
		returnFileName = r.FormValue("pre_thumb")
	}
	slug := tools.GenRandSlug(12)
	series, err := strconv.Atoi(r.FormValue("series_id"))
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	episodes, err := strconv.Atoi(r.FormValue("count_episodes"))
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	companyId, err := strconv.Atoi(r.FormValue("company_id"))
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	a := domain.AnimeInsert{
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
		CompanyId:     tools.NewNullInt(companyId),
		HashTag:       tools.NewNullString(r.FormValue("hash_tag")),
		TwitterUrl:    tools.NewNullString(r.FormValue("twitter_url")),
		OfficialUrl:   tools.NewNullString(r.FormValue("official_url")),
	}
	insertedId, err = adc.interactor.AnimeInsert(a)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": insertedId})
	return
}

func (adc *AdminController) AnimeUpdateView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("thumb")
	var returnFileName string
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = adc.s3.Image(file, fileHeader.Filename, []string{"anime"})
		if err != nil {
			response(w, r, perr.Wrap(err, perr.BadRequest), nil)
			return
		}
	} else {
		returnFileName = r.FormValue("pre_thumb")
	}
	// @TODO confirm OK?
	series, err := strconv.Atoi(r.FormValue("series_id"))
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	episodes, err := strconv.Atoi(r.FormValue("count_episodes"))
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	companyId, err := strconv.Atoi(r.FormValue("company_id"))
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	a := domain.AnimeInsert{
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
		CompanyId:     tools.NewNullInt(companyId),
		HashTag:       tools.NewNullString(r.FormValue("hash_tag")),
		TwitterUrl:    tools.NewNullString(r.FormValue("twitter_url")),
		OfficialUrl:   tools.NewNullString(r.FormValue("official_url")),
	}

	rowsAffected, err := adc.interactor.AnimeUpdate(id, a)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": rowsAffected})
	return
}

func (adc *AdminController) AnimeDeleteView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	rowsAffected, err := adc.interactor.AnimeDelete(id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": rowsAffected})
	return
}

/************************
         platform
*************************/

func (adc *AdminController) PlatformView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, err := strconv.Atoi(id)
		if err != nil {
			response(w, r, perr.Wrap(err, perr.BadRequest), nil)
			return
		}
		platform, err := adc.interactor.PlatformDetail(i)
		if err != nil {
			err = perr.Wrap(err, perr.NotFound)
		} else if platform.ID == 0 {
			err = perr.New("", perr.NotFound)
		}
		response(w, r, err, map[string]interface{}{"data": platform})
	} else {
		platforms, err := adc.interactor.PlatformAllAdmin()
		response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": platforms})
	}
	return
}

func (adc *AdminController) PlatformInsertView(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("image")
	var returnFileName string
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = adc.s3.Image(file, fileHeader.Filename, []string{"platform"})

		if err != nil {
			response(w, r, perr.Wrap(err, perr.BadRequest), nil)
			return
		}
	} else {
		returnFileName = ""
	}
	validStr := r.FormValue("valid")
	isValid, err := strconv.ParseBool(validStr)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	p := domain.TPlatform{
		EngName:  r.FormValue("engName"),
		PlatName: tools.NewNullString(r.FormValue("platName")),
		BaseUrl:  tools.NewNullString(r.FormValue("baseUrl")),
		Image:    tools.NewNullString(returnFileName),
		IsValid:  isValid,
	}
	lastInserted, err := adc.interactor.PlatformInsert(p)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": lastInserted})
	return
}

func (adc *AdminController) PlatformUpdateView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("image")
	var returnFileName string
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = adc.s3.Image(file, fileHeader.Filename, []string{"platform"})

		if err != nil {
			fmt.Print(err)
		}
	} else {
		returnFileName = ""
	}
	validStr := r.FormValue("valid")
	isValid, err := strconv.ParseBool(validStr)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	p := domain.TPlatform{
		EngName:  r.FormValue("engName"),
		PlatName: tools.NewNullString(r.FormValue("platName")),
		BaseUrl:  tools.NewNullString(r.FormValue("baseUrl")),
		Image:    tools.NewNullString(returnFileName),
		IsValid:  isValid,
	}
	rowsAffected, err := adc.interactor.PlatformUpdate(p, id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": rowsAffected})
	return
}

func (adc *AdminController) PlatformDeleteview(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	rowsAffected, err := adc.interactor.PlatformDelete(id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": rowsAffected})
	return
}

func (adc *AdminController) RelationPlatformView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id") // animeId
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	relations, err := adc.interactor.RelationPlatformByAnime(id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": relations})
	return
}

func (adc *AdminController) InsertRelationPlatformView(w http.ResponseWriter, r *http.Request) {
	var p domain.TRelationPlatformInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := adc.interactor.RelationPlatformInsert(p)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": lastInserted})
	return
}

func (adc *AdminController) DeleteRelationPlatformView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strAnimeId := query.Get("anime")
	strPlatformId := query.Get("platform")
	animeId, err := strconv.Atoi(strAnimeId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	platformId, err := strconv.Atoi(strPlatformId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	rowsAffected, err := adc.interactor.RelationPlatformDelete(animeId, platformId)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": rowsAffected})
	return
}

/************************
         season
*************************/

func (adc *AdminController) SeasonView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, err := strconv.Atoi(id)
		if err != nil {
			response(w, r, perr.Wrap(err, perr.BadRequest), nil)
			return
		}
		s, err := adc.interactor.DetailSeason(i)
		response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": s})
	} else {
		s, err := adc.interactor.ListSeason()
		response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": s})
	}
	return
}

func (adc *AdminController) InsertSeasonView(w http.ResponseWriter, r *http.Request) {
	var p domain.TSeasonInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := adc.interactor.InsertSeason(p)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": lastInserted})
	return
}

func (adc *AdminController) InsertRelationSeasonView(w http.ResponseWriter, r *http.Request) {
	var s domain.TSeasonRelationInput
	json.NewDecoder(r.Body).Decode(&s)
	lastInserted, err := adc.interactor.InsertRelationSeasonAnime(s)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": lastInserted})
	return
}

func (adc *AdminController) DeleteRelationSeasonView(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	rawAnimeId := qs.Get("anime")
	rawSeasonId := qs.Get("season")
	animeId, err := strconv.Atoi(rawAnimeId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	seasonId, err := strconv.Atoi(rawSeasonId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	affected, err := adc.interactor.DeleteRelationSeasonAnime(animeId, seasonId)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": affected})
	return
}

/************************
         series
*************************/

func (adc *AdminController) SeriesView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, err := strconv.Atoi(id)
		if err != nil {
			response(w, r, perr.Wrap(err, perr.BadRequest), nil)
			return
		}
		s, err := adc.interactor.DetailSeries(i)
		response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": s})
	} else {
		s, err := adc.interactor.ListSeries()
		response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": s})
	}
	return
}

func (adc *AdminController) InsertSeriesView(w http.ResponseWriter, r *http.Request) {
	var p domain.TSeriesInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := adc.interactor.InsertSeries(p)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": lastInserted})
	return
}

func (adc *AdminController) UpdateSeriesView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strId := query.Get("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}

	var p domain.TSeriesInput
	json.NewDecoder(r.Body).Decode(&p)
	rowsAffected, err := adc.interactor.UpdateSeries(p, id)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": rowsAffected})
	return
}

/************************
         company
*************************/

func (adc *AdminController) AdminCompanyView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	eng := query.Get("company")

	if eng != "" {
		company, err := adc.interactor.DetailCompany(eng)
		response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": company})
	} else {
		companies, err := adc.interactor.ListCompany()
		response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": companies})
	}
	return
}

func (adc *AdminController) InsertCompanyView(w http.ResponseWriter, r *http.Request) {
	var p domain.CompanyInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := adc.interactor.InsertCompany(p)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": lastInserted})
	return
}

func (adc *AdminController) UpdateCompanyView(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	eng := query.Get("company")

	var p domain.CompanyInput
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response(w, r, perr.Wrap(err, perr.BadRequest), nil)
		return
	}
	affected, err := adc.interactor.UpdateCompany(p, eng)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": affected})
	return
}

// @TODO:ADD delete

/************************
         staff
*************************/

func (adc *AdminController) InsertStaffView(w http.ResponseWriter, r *http.Request) {
	var p domain.StaffInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := adc.interactor.InsertStaff(p)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": lastInserted})
	return
}

// @TODO:ADD update and delete

/************************
         role
*************************/

func (adc *AdminController) InsertRoleView(w http.ResponseWriter, r *http.Request) {
	var p domain.RoleInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := adc.interactor.InsertRole(p)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": lastInserted})
	return
}

func (adc *AdminController) ListRoleView(w http.ResponseWriter, r *http.Request) {
	roles, err := adc.interactor.RoleList()
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": roles})
	return
}

func (adc *AdminController) InsertStaffRoleView(w http.ResponseWriter, r *http.Request) {
	var p domain.AnimeStaffRoleInput
	json.NewDecoder(r.Body).Decode(&p)
	lastInserted, err := adc.interactor.InsertStaffRole(p)
	response(w, r, perr.Wrap(err, perr.BadRequest), map[string]interface{}{"data": lastInserted})
	return
}

// @TODO:ADD update and delete
