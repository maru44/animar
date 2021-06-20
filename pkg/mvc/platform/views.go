package platform

import (
	"animar/v1/pkg/tools/api"
	"animar/v1/pkg/tools/s3"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func PlatformView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	id := query.Get("id")

	if id != "" {
		i, _ := strconv.Atoi(id)
		plat := detailPlatfrom(i)
		if plat.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return errors.New("Not Found")
		}
		api.JsonResponse(w, map[string]interface{}{"data": plat})
	} else {
		plats := ListPlatformDomain()
		api.JsonResponse(w, map[string]interface{}{"data": plats})
	}
	return nil
}

func InsertPlatformView(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB

	file, fileHeader, err := r.FormFile("image")
	var returnFileName string
	if err == nil {
		// w/ thumb picture
		defer file.Close()
		returnFileName, err = s3.UploadS3(file, fileHeader.Filename, []string{"platform"})

		if err != nil {
			fmt.Print(err)
		}
	} else {
		returnFileName = ""
	}
	validStr := r.FormValue("valid")
	isValid, _ := strconv.ParseBool(validStr)
	insertedId := insertPlatform(
		r.FormValue("engName"), r.FormValue("platName"), r.FormValue("baseUrl"),
		returnFileName, isValid,
	)
	api.JsonResponse(w, map[string]interface{}{"data": insertedId})
	return nil
}

// update ?id=<id>
func UpdatePlatformView(w http.ResponseWriter, r *http.Request) error {
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
		returnFileName, err = s3.UploadS3(file, fileHeader.Filename, []string{"platform"})

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
	api.JsonResponse(w, map[string]interface{}{"data": updatedId})
	return nil
}

// delete platform ?=<id>
func DeletePlatformView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strId := query.Get("id")
	id, _ := strconv.Atoi(strId)

	deletedRow := deletePlatform(id)
	api.JsonResponse(w, map[string]interface{}{"data": deletedRow})
	return nil
}

/****************************
*          relation		    *
****************************/

func RelationPlatformView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	id := query.Get("id") // animeId
	i, _ := strconv.Atoi(id)
	relations := ListRelationPlatformDomain(i)
	api.JsonResponse(w, map[string]interface{}{"data": relations})
	return nil
}

func InsertRelationPlatformView(w http.ResponseWriter, r *http.Request) error {
	var p TRelationPlatformInput
	json.NewDecoder(r.Body).Decode(&p)
	value := insertRelation(
		p.PlatformId, p.AnimeId, p.LinkUrl,
	)
	api.JsonResponse(w, map[string]interface{}{"data": value})
	return nil
}

// delete platform ?=<id>
func DeleteRelationPlatformView(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	strAnimeId := query.Get("anime")
	strPlatformId := query.Get("platform")
	animeId, _ := strconv.Atoi(strAnimeId)
	platformId, _ := strconv.Atoi(strPlatformId)

	deletedRow := deleteRelationPlatform(animeId, platformId)
	api.JsonResponse(w, map[string]interface{}{"data": deletedRow})
	return nil
}
