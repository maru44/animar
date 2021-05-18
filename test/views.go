package test

import (
	"animar/v1/helper"
	"encoding/json"
	"fmt"
	"net/http"
)

type TInputFile struct {
	File     string `json:"file"`
	FileName string `json:"filename"`
}

func Uploader(w http.ResponseWriter, r *http.Request) error {
	result := helper.TVoidJsonResponse{Status: 200}
	var posted TInputFile
	json.NewDecoder(r.Body).Decode(&posted)

	fileName, err := helper.UploadS3([]byte(posted.File), posted.FileName, "test")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(fileName)

	result.ResponseWrite(w)
	return nil
}
