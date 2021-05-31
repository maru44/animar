package blog

import (
	"animar/v1/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

type TBlogInputWith struct {
	Title    string `json:"Title"`
	Abstract string `json:"Abstract,omitempty"`
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
	tools.SetDefaultResponseHeader(w)
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
	tools.SetDefaultResponseHeader(w)
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
	tools.SetDefaultResponseHeader(w)
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
		blogs = ListBlogByUserJoinAnimeUserDomain(uid)
	} else {
		blogs = ListBlogJoinAnimeUserDomain()
	}

	result.Data = blogs
	result.ResponseWrite(w)
	return nil
}

func InsertBlogView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetIdFromCookie(r)

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
	result := tools.TIntJsonReponse{Status: 200}
	userId := tools.GetIdFromCookie(r)

	if userId == "" {
		result.Status = 4001
	} else {
		var posted TBlogInputWith
		json.NewDecoder(r.Body).Decode(&posted)
		// @TODO delete
		fmt.Print(posted)

		value := InsertBlog(posted.Title, posted.Abstract, posted.Content, userId)
		result.Num = value

		for _, animeId := range posted.AnimeIds {
			InsertRelationAnimeBlog(animeId, value)
		}
	}
	result.ResponseWrite(w)
	return nil
}

func UpdateBlogWithRelationView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TIntJsonReponse{Status: 200}
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

	var posted TBlogInputWith
	json.NewDecoder(r.Body).Decode(&posted)
	value := UpdateBlog(id, posted.Title, posted.Abstract, posted.Content)
	UpdateRelationBlogAnimesDomain(posted.AnimeIds, id)
	result.Num = value

	result.ResponseWrite(w)
	return nil
}

func DeleteBlogView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TIntJsonReponse{Status: 200}
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
	result.Num = deletedRow
	result.ResponseWrite(w)
	return nil
}
