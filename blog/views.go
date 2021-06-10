package blog

import (
	"animar/v1/tools"
	"encoding/json"
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
	result := tools.TBaseJsonResponse{Status: 200}

	var blogs []TBlog
	blogs = ListBlogDomain()

	result.Data = blogs
	result.ResponseWrite(w)
	return nil
}

// retrieve blog + anime
func BlogJoinAnimeView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	// access user
	var userId string
	switch r.Method {
	case "GET":
		userId = tools.GetIdFromCookie(r)
	case "POST":
		var posted tools.TUserIdCookieInput
		json.NewDecoder(r.Body).Decode(&posted)
		userId = tools.GetUserIdFromToken(posted.Token)
	default:
		userId = ""
	}

	query := r.URL.Query()
	slug := query.Get("s")
	id := query.Get("id")
	uid := query.Get("u")

	var blogs []TBlogJoinAnimes
	if slug != "" {
		blog := DetailBlogJoinAnimeDomain(slug)
		if blog.ID == 0 {
			result.Status = 404
		}
		blogs = append(blogs, blog)
	} else if id != "" {
		i, _ := strconv.Atoi(id)
		blog := DetailBlogJoinAnimeFromIdDomain(i)
		if blog.ID == 0 {
			result.Status = 404
		}
		blogs = append(blogs, blog)
	} else if uid != "" {
		blogs = ListBlogByUserJoinAnimeDomain(uid, userId)
	} else {
		blogs = ListBlogJoinAnimeDomain()
	}

	result.Data = blogs
	result.ResponseWrite(w)
	return nil
}

// retrieve blog + anime + user
func BlogJoinAnimeUserView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	// access user
	var userId string
	switch r.Method {
	case "GET":
		userId = tools.GetIdFromCookie(r)
	case "POST":
		var posted tools.TUserIdCookieInput
		json.NewDecoder(r.Body).Decode(&posted)
		userId = tools.GetUserIdFromToken(posted.Token)
	default:
		userId = ""
	}

	query := r.URL.Query()
	slug := query.Get("s")
	id := query.Get("id")
	uid := query.Get("u")

	var blogs []TBlogJoinAnimesUser
	if slug != "" {
		blog := DetailBlogJoinAnimeUserDomain(slug)
		if blog.ID == 0 {
			result.Status = 404
		}
		blogs = append(blogs, blog)
	} else if id != "" {
		i, _ := strconv.Atoi(id)
		blog := DetailBlogJoinAnimeUserFromIdDomain(i)
		if blog.ID == 0 {
			result.Status = 404
		}
		blogs = append(blogs, blog)
	} else if uid != "" {
		blogs = ListBlogByUserJoinAnimeUserDomain(uid, userId)
	} else {
		blogs = ListBlogJoinAnimeUserDomain()
	}

	result.Data = blogs
	result.ResponseWrite(w)
	return nil
}

func InsertBlogWithRelationView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	userId := tools.GetIdFromCookie(r)

	if userId == "" {
		result.Status = 4001
	} else {
		var p TBlogInputWith
		json.NewDecoder(r.Body).Decode(&p)

		value := InsertBlog(p.Title, p.Abstract, p.Content, userId, p.IsPublic)
		result.Data = value

		for _, animeId := range p.AnimeIds {
			InsertRelationAnimeBlog(animeId, value)
		}
	}
	result.ResponseWrite(w)
	return nil
}

func UpdateBlogWithRelationView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	userId := tools.GetIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId := BlogUserId(id)
	if blogUserId != userId {
		result.Status = 4003
		result.ResponseWrite(w)
		return nil
	}

	var p TBlogInputWith
	json.NewDecoder(r.Body).Decode(&p)
	value := UpdateBlog(id, p.Title, p.Abstract, p.Content, p.IsPublic)
	UpdateRelationBlogAnimesDomain(p.AnimeIds, id)
	result.Data = value

	result.ResponseWrite(w)
	return nil
}

func DeleteBlogView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	userId := tools.GetIdFromCookie(r)

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	// user 不一致
	blogUserId := BlogUserId(id)
	if blogUserId != userId {
		result.Status = 4003
		result.ResponseWrite(w)
		return nil
	}

	deletedRow := DeleteBlog(id)
	result.Data = deletedRow
	result.ResponseWrite(w)
	return nil
}
