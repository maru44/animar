package anime

import (
	"animar/v1/tools/api"
	"animar/v1/tools/s3"
	"animar/v1/tools/tools"
	"encoding/json"
	"errors"
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

	query := r.URL.Query()
	strId := query.Get("id")
	slug := query.Get("slug")
	season := query.Get("season")

	if strId != "" {
		id, _ := strconv.Atoi(strId)
		ani := DetailAnime(id)
		if ani.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": ani})
	} else if slug != "" {
		ani := DetailAnimeBySlug(slug)
		if ani.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": ani})
	} else if season != "" {
		seasonId, _ := strconv.Atoi(season)
		animes := AnimesBySeasonIdDomain(seasonId)
		api.JsonResponse(w, map[string]interface{}{"data": animes})
	} else {
		animes := ListAnimeDomain()
		api.JsonResponse(w, map[string]interface{}{"data": animes})
	}
	switch {
	case strId != "":
		id, _ := strconv.Atoi(strId)
		ani := DetailAnime(id)
		if ani.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": ani})
	case slug != "":
		ani := DetailAnimeBySlug(slug)
		if ani.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": ani})
	}
	return nil
}

// not be used to
// anime data + user's watch + review
// func AnimeWithUserWatchView(w http.ResponseWriter, r *http.Request) error {
// 	result := api.TBaseJsonResponse{Status: 200}

// 	query := r.URL.Query()
// 	strId := query.Get("id")
// 	slug := query.Get("slug")
// 	userId := fire.GetIdFromCookie(r)

// 	var animes []TAnimeWithUserWatchReview
// 	if strId != "" {
// 		/*
// 			id, _ := strconv.Atoi(strId)
// 			ani := DetailAnime(id)
// 			if ani.ID == 0 {
// 				result.Status = 404
// 			}
// 			animes = append(animes, ani)
// 		*/
// 	} else if slug != "" {
// 		ani := DetailAnimeBySlugWithUserWatchReview(slug, userId)
// 		if ani.ID == 0 {
// 			w.WriteHeader(http.StatusNotFound)
// 			return errors.New("Not Found")
// 		} else {
// 			animes = append(animes, ani)
// 			result.Data = animes
// 		}
// 	} else {
// 		// animes = ListAnimeDomain()
// 	}

// 	result.ResponseWrite(w)
// 	return nil
// }

func ListAnimeMinimumView(w http.ResponseWriter, r *http.Request) error {
	var animes []TAnimeMinimum
	animes = ListAnimeMinimumDomain()

	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}

func SearchAnimeMinView(w http.ResponseWriter, r *http.Request) error {
	var animes []TAnimeMinimum

	query := r.URL.Query()
	title := query.Get("t")

	animes = ListAnimeMinimumDomainByTitle(title)
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}

/************************************
             for admin
************************************/

type TUserToken struct {
	Token string `json:"token,omitempty"`
}

func AnimeListAdminView(w http.ResponseWriter, r *http.Request) error {
	var animes []TAnime
	animes = ListAnimeDomain()
	api.JsonResponse(w, map[string]interface{}{"data": animes})
	return nil
}

// anime detail ?=<id>
func AnimeDetailAdminView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	var posted TUserToken
	json.NewDecoder(r.Body).Decode(&posted)

	ani := DetailAnimeAdmin(id)
	if ani.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return errors.New("Not Found")
	}
	api.JsonResponse(w, map[string]interface{}{"data": ani})
	return nil
}

// add anime (only admin)
func AnimePostView(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	var returnFileName string
	var insertedId int
	file, fileHeader, err := r.FormFile("thumb")
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
	insertedId = InsertAnime(
		r.FormValue("title"), r.FormValue("abbreviation"),
		r.FormValue("kana"), r.FormValue("engName"),
		r.FormValue("description"), r.FormValue("state"),
		series, episodes, r.FormValue("copyright"), returnFileName,
	)
	api.JsonResponse(w, map[string]interface{}{"data": insertedId})
	return nil
}

// update ?=<id>
func AnimeUpdateView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("thumb")
	var returnFileName string
	var updatedId int

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
	updatedId = UpdateAnime(
		id, r.FormValue("title"), r.FormValue("abbreviation"),
		r.FormValue("kana"), r.FormValue("engName"),
		r.FormValue("description"), r.FormValue("state"),
		series, episodes, r.FormValue("copyright"), returnFileName,
	)
	api.JsonResponse(w, map[string]interface{}{"data": updatedId})
	return nil
}

// delete anime ?=<id>
func AnimeDeleteView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	deletedRow := DeleteAnime(id)
	api.JsonResponse(w, map[string]interface{}{"data": deletedRow})
	return nil
}
