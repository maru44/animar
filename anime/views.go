package anime

import (
	"animar/v1/tools"
	"encoding/json"
	"net/http"
	"strconv"
)

type TAnimeInput struct {
	Title         string `json:"title"`
	Abbrevation   string `json:"abbrevation,omitempty"`
	Kana          string `json:"kana,omitempty"`
	EngName       string `json:"eng_name:omitempty"`
	ThumbUrl      string `json:"thumb_url,omitempty"`
	PreThumbUrl   string `json:"pre_thumb,omitempty"`
	Description   string `json:"description,omitempty"`
	State         int    `json:"state,omitempty"`
	SeriesId      int    `json:"series_id,omitempty"`
	CountEpisodes int    `jsoin:"count_episodes,omitempty"`
}

// list and detail
func AnimeView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	query := r.URL.Query()
	strId := query.Get("id")
	slug := query.Get("slug")

	var animes []TAnime
	if strId != "" {
		id, _ := strconv.Atoi(strId)
		ani := DetailAnime(id)
		if ani.ID == 0 {
			result.Status = 404
		}
		animes = append(animes, ani)
		result.Data = animes
	} else if slug != "" {
		var animes []TAnimeWithSeries
		ani := DetailAnimeBySlug(slug)
		if ani.ID == 0 {
			result.Status = 404
		}
		animes = append(animes, ani)
		result.Data = animes
	} else {
		animes = ListAnimeDomain()
		result.Data = animes
	}

	result.ResponseWrite(w)

	return nil
}

// not be used to
// anime data + user's watch + review
func AnimeWithUserWatchView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	query := r.URL.Query()
	strId := query.Get("id")
	slug := query.Get("slug")
	userId := tools.GetIdFromCookie(r)

	var animes []TAnimeWithUserWatchReview
	if strId != "" {
		/*
			id, _ := strconv.Atoi(strId)
			ani := DetailAnime(id)
			if ani.ID == 0 {
				result.Status = 404
			}
			animes = append(animes, ani)
		*/
	} else if slug != "" {
		ani := DetailAnimeBySlugWithUserWatchReview(slug, userId)
		if ani.ID == 0 {
			result.Status = 404
		} else {
			animes = append(animes, ani)
			result.Data = animes
		}
	} else {
		// animes = ListAnimeDomain()
	}

	result.ResponseWrite(w)
	return nil
}

func ListAnimeMinimumView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	var animes []TAnimeMinimum

	animes = ListAnimeMinimumDomain()

	result.Data = animes
	result.ResponseWrite(w)
	return nil
}

func SearchAnimeView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	var animes []TAnimeMinimum

	query := r.URL.Query()
	title := query.Get("t")

	animes = ListAnimeMinimumDomainByTitle(title)
	result.Data = animes
	result.ResponseWrite(w)
	return nil
}

/************************************
             for admin
************************************/

type TUserToken struct {
	Token string `json:"token,omitempty"`
}

func AnimeListAdminView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	var userId string
	switch r.Method {
	case "GET":
		userId = tools.GetAdminIdFromCookie(r)
	case "POST":
		var posted TUserToken
		json.NewDecoder(r.Body).Decode(&posted)
		userId = tools.GetAdminIdFromIdToken(posted.Token)
	default:
		userId = ""
	}

	var animes []TAnime
	if userId == "" {
		result.Status = 4003
	} else {
		animes = ListAnimeDomain()
	}
	result.Data = animes
	result.ResponseWrite(w)
	return nil
}

// anime detail ?=<id>
func AnimeDetailAdminView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	var userId string
	switch r.Method {
	case "GET":
		userId = tools.GetAdminIdFromCookie(r)
	case "POST":
		var posted TUserToken
		json.NewDecoder(r.Body).Decode(&posted)
		userId = tools.GetAdminIdFromIdToken(posted.Token)
	default:
		userId = ""
	}

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	var posted TUserToken
	json.NewDecoder(r.Body).Decode(&posted)

	var animes []TAnimeAdmin
	if userId == "" {
		result.Status = 4003
	} else {
		ani := DetailAnimeAdmin(id)
		if ani.ID == 0 {
			result.Status = 404
		} else {
			animes = append(animes, ani)
		}
	}
	result.Data = animes
	result.ResponseWrite(w)
	return nil
}

// add anime (only admin)
func AnimePostView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	result, is_valid := result.LimitMethod([]string{"POST"}, r)
	if !is_valid {
		result.ResponseWrite(w)
		return nil
	}

	userId := tools.GetAdminIdFromCookie(r)
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	if userId == "" {
		result.Status = 4003
	} else {
		file, fileHeader, err := r.FormFile("thumb")
		var returnFileName string
		var insertedId int
		if err == nil {
			// w/ thumb picture
			defer file.Close()
			returnFileName, err = tools.UploadS3(file, fileHeader.Filename, []string{"anime"})

			if err != nil {
				tools.ErrorLog(err)
			}
		} else {
			returnFileName = ""
		}
		series, _ := strconv.Atoi(r.FormValue("series_id"))
		episodes, _ := strconv.Atoi(r.FormValue("count_episodes"))
		insertedId = InsertAnime(
			r.FormValue("title"), r.FormValue("abbreviation"),
			r.FormValue("kana"), r.FormValue("engName"),
			r.FormValue("description"), r.FormValue("state"),
			series, episodes, r.FormValue("copyright"), returnFileName,
		)
		result.Data = insertedId
	}
	result.ResponseWrite(w)
	return nil
}

// update ?=<id>
func AnimeUpdateView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	result, is_valid := result.LimitMethod([]string{"PUT"}, r)
	if !is_valid {
		result.ResponseWrite(w)
		return nil
	}

	userId := tools.GetAdminIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	if userId == "" {
		result.Status = 4003
	} else {
		file, fileHeader, err := r.FormFile("thumb")
		var returnFileName string
		var updatedId int

		if err == nil {
			// w/ thumb picture
			defer file.Close()
			returnFileName, err = tools.UploadS3(file, fileHeader.Filename, []string{"anime"})
			if err != nil {
				tools.ErrorLog(err)
			}
		} else {
			returnFileName = r.FormValue("pre_thumb")
		}

		series, _ := strconv.Atoi(r.FormValue("series_id"))
		episodes, _ := strconv.Atoi(r.FormValue("count_episodes"))
		updatedId = UpdateAnime(
			id, r.FormValue("title"), r.FormValue("abbreviation"),
			r.FormValue("kana"), r.FormValue("engName"),
			r.FormValue("description"), r.FormValue("state"),
			series, episodes, r.FormValue("copyright"), returnFileName,
		)
		result.Data = updatedId
	}
	result.ResponseWrite(w)
	return nil
}

// delete anime ?=<id>
func AnimeDeleteView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	result, is_valid := result.LimitMethod([]string{"DELETE"}, r)
	if !is_valid {
		result.ResponseWrite(w)
		return nil
	}

	userId := tools.GetAdminIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	if userId == "" {
		result.Status = 4003
	} else {
		deletedRow := DeleteAnime(id)
		result.Data = deletedRow
	}

	result.ResponseWrite(w)
	return nil
}
