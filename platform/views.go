package platform

import (
	"animar/v1/tools"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func PlatformView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	query := r.URL.Query()
	id := query.Get("id")

	var plats []TPlatform
	if id != "" {
		i, _ := strconv.Atoi(id)
		plat := detailPlatfrom(i)
		if plat.ID == 0 {
			result.Status = 404
		}
		plats = append(plats, plat)
	} else {
		plats = ListPlatformDomain()
	}
	result.Data = plats

	result.ResponseWrite(w)
	return nil
}

func InsertPlatformView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("image")
	var returnFileName string
	var insertedId int
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = tools.UploadS3(file, fileHeader.Filename, []string{"platform"})

		if err != nil {
			fmt.Print(err)
		}
	} else {
		returnFileName = ""
	}
	validStr := r.FormValue("valid")
	isValid, _ := strconv.ParseBool(validStr)
	insertedId = insertPlatform(
		r.FormValue("engName"), r.FormValue("platName"), r.FormValue("baseUrl"),
		returnFileName, isValid,
	)
	result.Data = insertedId

	result.ResponseWrite(w)
	return nil
}

// update ?id=<id>
func UpdatePlatformView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("image")
	var returnFileName string
	var updatedId int
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = tools.UploadS3(file, fileHeader.Filename, []string{"platform"})

		if err != nil {
			fmt.Print(err)
		}
	} else {
		returnFileName = ""
	}
	validStr := r.FormValue("valid")
	isValid, _ := strconv.ParseBool(validStr)
	updatedId = updatePlatform(
		r.FormValue("engName"), r.FormValue("platName"), r.FormValue("baseUrl"),
		returnFileName, isValid, id,
	)
	result.Data = updatedId

	result.ResponseWrite(w)
	return nil
}

// delete platform ?=<id>
func DeletePlatformView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	deletedRow := deletePlatform(id)
	result.Data = deletedRow

	result.ResponseWrite(w)
	return nil
}

/****************************
*          relation		    *
****************************/

func RelationPlatformView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}
	query := r.URL.Query()
	id := query.Get("id") // animeId
	i, _ := strconv.Atoi(id)
	relations := ListRelationPlatformDomain(i)
	result.Data = relations
	result.ResponseWrite(w)
	return nil
}

func InsertRelationPlatformView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	var p TRelationPlatformInput
	json.NewDecoder(r.Body).Decode(&p)
	value := insertRelation(
		p.PlatformId, p.AnimeId, p.LinkUrl,
	)
	result.Data = value

	result.ResponseWrite(w)
	return nil
}

// delete platform ?=<id>
func DeleteRelationPlatformView(w http.ResponseWriter, r *http.Request) error {
	result := tools.TBaseJsonResponse{Status: 200}

	query := r.URL.Query()
	strAnimeId := query.Get("anime")
	strPlatformId := query.Get("platform")
	animeId, _ := strconv.Atoi(strAnimeId)
	platformId, _ := strconv.Atoi(strPlatformId)

	deletedRow := deleteRelationPlatform(animeId, platformId)
	result.Data = deletedRow

	result.ResponseWrite(w)
	return nil
}
