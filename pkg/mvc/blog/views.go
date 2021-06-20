package blog

import (
	"animar/v1/tools/api"
	"animar/v1/tools/fire"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type TBlogInput struct {
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Content  string `json:"content"`
}

type TBlogInputWith struct {
	Title    string `json:"title"`
	Abstract string `json:"abstract,omitempty"`
	Content  string `json:"content"`
	IsPublic bool   `json:"is_public,omitempty"`
	AnimeIds []int  `json:"anime_ids"`
}

func ListBlogView(w http.ResponseWriter, r *http.Request) error {
	var blogs []TBlog
	blogs = ListBlogDomain()
	api.JsonResponse(w, map[string]interface{}{"data": blogs})
	return nil
}

// retrieve blog + anime
func BlogJoinAnimeView(w http.ResponseWriter, r *http.Request) error {
	// access user
	var userId string
	switch r.Method {
	case "GET":
		userId = fire.GetIdFromCookie(r)
	case "POST":
		var posted fire.TUserIdCookieInput
		json.NewDecoder(r.Body).Decode(&posted)
		userId = fire.GetUserIdFromToken(posted.Token)
	default:
		userId = ""
	}

	query := r.URL.Query()
	slug := query.Get("s")
	id := query.Get("id")
	uid := query.Get("u")

	if slug != "" {
		blog := DetailBlogJoinAnimeDomain(slug)
		if blog.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": blog})
	} else if id != "" {
		i, _ := strconv.Atoi(id)
		blog := DetailBlogJoinAnimeFromIdDomain(i)
		if blog.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": blog})
	} else if uid != "" {
		blogs := ListBlogByUserJoinAnimeDomain(uid, userId)
		api.JsonResponse(w, map[string]interface{}{"data": blogs})
	} else {
		blogs := ListBlogJoinAnimeDomain()
		api.JsonResponse(w, map[string]interface{}{"data": blogs})
	}
	return nil
}

// retrieve blog + anime + user
func BlogJoinAnimeUserView(w http.ResponseWriter, r *http.Request) error {
	// access user
	var userId string
	switch r.Method {
	case "GET":
		userId = fire.GetIdFromCookie(r)
	case "POST":
		var posted fire.TUserIdCookieInput
		json.NewDecoder(r.Body).Decode(&posted)
		userId = fire.GetUserIdFromToken(posted.Token)
	default:
		userId = ""
	}

	query := r.URL.Query()
	slug := query.Get("s")
	id := query.Get("id")
	uid := query.Get("u")

	if slug != "" {
		blog := DetailBlogJoinAnimeUserDomain(slug)
		if blog.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": blog})
	} else if id != "" {
		i, _ := strconv.Atoi(id)
		blog := DetailBlogJoinAnimeUserFromIdDomain(i)
		if blog.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": blog})
	} else if uid != "" {
		blogs := ListBlogByUserJoinAnimeUserDomain(uid, userId)
		api.JsonResponse(w, map[string]interface{}{"data": blogs})
	} else {
		blogs := ListBlogJoinAnimeUserDomain()
		api.JsonResponse(w, map[string]interface{}{"data": blogs})
	}
	return nil
}

func InsertBlogWithRelationView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return errors.New("Unauthorized")
	} else {
		var p TBlogInputWith
		json.NewDecoder(r.Body).Decode(&p)
		value := InsertBlog(p.Title, p.Abstract, p.Content, userId, p.IsPublic)
		for _, animeId := range p.AnimeIds {
			InsertRelationAnimeBlog(animeId, value)
		}
		api.JsonResponse(w, map[string]interface{}{"data": value})
	}
	return nil
}

func UpdateBlogWithRelationView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId := BlogUserId(id)
	if blogUserId != userId {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("Forbidden")
	}

	var p TBlogInputWith
	json.NewDecoder(r.Body).Decode(&p)
	value := UpdateBlog(id, p.Title, p.Abstract, p.Content, p.IsPublic)
	UpdateRelationBlogAnimesDomain(p.AnimeIds, id)

	api.JsonResponse(w, map[string]interface{}{"data": value})
	return nil
}

func DeleteBlogView(w http.ResponseWriter, r *http.Request) error {
	userId := fire.GetIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId := BlogUserId(id)
	if blogUserId != userId {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("Forbidden")
	}

	deletedRow := DeleteBlog(id)
	api.JsonResponse(w, map[string]interface{}{"data": deletedRow})
	return nil
}
