package blog

import (
	"animar/v1/helper"
	"encoding/json"
	"net/http"
)

type TBlogResponse struct {
	Status int     `json:"Status"`
	Data   []TBlog `json:"Data"`
}

type TBlogJoinAnimesResponse struct {
	Status int               `json:"Status"`
	Data   []TBlogJoinAnimes `json:"Data"`
}

type TBlogJoinAnimesUserResponse struct {
	Status int                   `json:"Status"`
	Data   []TBlogJoinAnimesUser `json:"Data"`
}

type TBlogInput struct {
	Title    string `json:"Title"`
	Abstract string `json:"Abstract"`
	Content  string `json:"Content"`
}

type TBlogInputWithResponse struct {
	Title    string `json:"Title"`
	Abstract string `json:"Abstract"`
	Content  string `json:"Content"`
	AnimeIds []int  `json:"anime_ids"`
}

// @TODO
// retrieve(list)とretrieve(detail)で型を分けて
// select して最適化

func (result TBlogResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	helper.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func (result TBlogJoinAnimesResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	helper.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func (result TBlogJoinAnimesUserResponse) ResponseWrite(w http.ResponseWriter) bool {
	res, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	helper.SetDefaultResponseHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	return true
}

func ListBlogView(w http.ResponseWriter, r *http.Request) error {
	result := TBlogResponse{Status: 200}

	var blogs []TBlog
	blogs = ListBlogDomain()

	result.Data = blogs
	result.ResponseWrite(w)
	return nil
}

func BlogJoinAnimeView(w http.ResponseWriter, r *http.Request) error {
	result := TBlogJoinAnimesUserResponse{Status: 200}

	query := r.URL.Query()
	slug := query.Get("s")
	uid := query.Get("u")

	var blogs []TBlogJoinAnimesUser
	if slug != "" {
		blog := DetailBlogJoinAnimeUserDomain(slug)
		if blog.ID == 0 {
			result.Status = 404
		}
		blogs = append(blogs, blog)
	} else if uid != "" {
		blogs = ListBlogByUserJoinAnimeUserDomain(uid)
	} else {
		blogs = ListBlogJoinAnimeUserDomain()
	}

	result.Data = blogs
	result.ResponseWrite(w)
	return nil
}

func InsertBlogView(w http.ResponseWriter, r *http.Request) error {
	result := helper.TIntJsonReponse{Status: 200}
	userId := helper.GetIdFromCookie(r)

	if userId == "" {
		result.Status = 4001
	} else {
		var posted TBlogInput
		json.NewDecoder(r.Body).Decode(&posted)
		value := InsertBlog(posted.Title, posted.Abstract, posted.Content, userId)
		result.Num = value
	}
	result.ResponseWrite(w)
	return nil
}

func InsertBlogWithRelationView(w http.ResponseWriter, r *http.Request) error {
	result := helper.TIntJsonReponse{Status: 200}
	userId := helper.GetIdFromCookie(r)

	if userId == "" {
		result.Status = 4001
	} else {
		var posted TBlogInputWithResponse
		json.NewDecoder(r.Body).Decode(&posted)
		value := InsertBlog(posted.Title, posted.Abstract, posted.Content, userId)
		result.Num = value

		for _, animeId := range posted.AnimeIds {
			InsertRelationAnimeBlog(animeId, value)
		}
	}
	result.ResponseWrite(w)
	return nil
}
