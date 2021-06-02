package anime

import (
	"animar/v1/tools"
	"encoding/json"
	"net/http"
	"strconv"
)

type TAnimesJsonResponse struct {
	Status int      `json:"Status"`
	Data   []TAnime `json:"Data"`
}

type TAnimesWithUserWatchResponse struct {
	Status int                         `json:"Status"`
	Data   []TAnimeWithUserWatchReview `json:"Data"`
}

type TAnimeInput struct {
	Title       string `json:"Title"`
	Abbrevation string `json:"Abbrevation,omitempty"`
	Kana        string `json:"Kana,omitempty"`
	EngName     string `json:"EngName:omitempty"`
	ThumbUrl    string `json:"ThumbUrl,omitempty"`
	Content     string `json:"Content,omitempty"`
	OnAirState  int    `json:"OnAirState,omitempty"`
	SeriesId    int    `json:"Series,omitempty"`
	Season      string `jsoin:"Season,omitempty"`
	Stories     int    `jsoin:"Stories,omitempty"`
}

type TAnimeMinimumResponse struct {
	Status int             `json:"Status"`
	Data   []TAnimeMinimum `json:"Data"`
}

type TAnimeAdminResponse struct {
	Status int           `json:"Status"`
	Data   []TAnimeAdmin `json:"Data"`
}

func (animeJson TAnimesJsonResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(animeJson)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	tools.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func (animeWUWCJson TAnimesWithUserWatchResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(animeWUWCJson)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	tools.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func (result TAnimeMinimumResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	tools.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func (result TAnimeAdminResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	tools.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

// list and detail
func AnimeView(w http.ResponseWriter, r *http.Request) error {
	result := TAnimesJsonResponse{Status: 200}

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
	} else if slug != "" {
		ani := DetailAnimeBySlug(slug)
		if ani.ID == 0 {
			result.Status = 404
		}
		animes = append(animes, ani)
	} else {
		animes = ListAnimeDomain()
	}

	result.Data = animes
	result.ResponseWrite(w)

	return nil
}

// not be used to
// anime data + user's watch + review
func AnimeWithUserWatchView(w http.ResponseWriter, r *http.Request) error {
	result := TAnimesWithUserWatchResponse{Status: 200}

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
	result := TAnimeMinimumResponse{Status: 200}
	var animes []TAnimeMinimum

	animes = ListAnimeMinimumDomain()

	result.Data = animes
	result.ResponseWrite(w)
	return nil
}

func SearchAnimeView(w http.ResponseWriter, r *http.Request) error {
	result := TAnimeMinimumResponse{Status: 200}
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

// anime detail ?=<id>
func AnimeDetailAdminView(w http.ResponseWriter, r *http.Request) error {
	result := TAnimeAdminResponse{Status: 200}
	userId := tools.GetAdminIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

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
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetAdminIdFromCookie(r)
	if userId == "" {
		result.Status = 4003
	} else {
		var posted TAnimeInput
		json.NewDecoder(r.Body).Decode(&posted)
		insertedId := InsertAnime(
			posted.Title, posted.Abbrevation, posted.Kana, posted.EngName, posted.Content, posted.OnAirState,
			posted.SeriesId, posted.Season, posted.Stories, posted.ThumbUrl,
		)
		result.Num = insertedId
	}
	result.ResponseWrite(w)
	return nil
}

// update ?=<id>
func AnimeUpdateView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetAdminIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	if userId == "" {
		result.Status = 4003
	} else {
		var posted TAnimeInput
		json.NewDecoder(r.Body).Decode(&posted)
		updatedId := UpdateAnime(
			id, posted.Abbrevation, posted.Kana, posted.EngName,
			posted.ThumbUrl, posted.Title, posted.Content, posted.OnAirState,
			posted.SeriesId, posted.Season, posted.Stories,
		)
		result.Num = updatedId
	}
	result.ResponseWrite(w)
	return nil
}

// delete anime ?=<id>
func AnimeDeleteView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetAdminIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	if userId == "" {
		result.Status = 4003
	} else {
		deletedRow := DeleteAnime(id)
		result.Num = deletedRow
	}

	result.ResponseWrite(w)
	return nil
}
