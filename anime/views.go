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
	Abbrevation string `json:"Abbrevation"`
	ThumbUrl    string `json:"ThumbUrl"`
	Content     string `json:"Content"`
	OnAirState  int    `json:"OnAirState"`
	SeriesId    int    `json:"Series"`
	Season      string `jsoin:"Season"`
	Stories     int    `jsoin:"Stories"`
}

type TAnimeMinimumResponse struct {
	Status int             `json:"Status"`
	Data   []TAnimeMinimum `json:"Data"`
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
		}
		animes = append(animes, ani)
	} else {
		// animes = ListAnimeDomain()
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
			posted.Title, posted.Abbrevation, posted.Content, posted.OnAirState,
			posted.SeriesId, posted.Season, posted.Stories, posted.ThumbUrl,
		)
		result.Num = insertedId
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
