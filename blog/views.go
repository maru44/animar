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

type TBlogInput struct {
	Title    string `json:"Title"`
	Abstract string `json:"Abstract"`
	Content  string `json:"Content"`
}

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

func ListBlogView(w http.ResponseWriter, r *http.Request) error {
	result := TBlogResponse{Status: 200}

	var blogs []TBlog
	blogs = ListBlogDomain()

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
