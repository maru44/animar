package test

import (
	"animar/v1/helper"
	"fmt"
	"net/http"
	"reflect"
)

type TInputFile struct {
	File     string `json:"file"`
	FileName string `json:"filename"`
}

func Uploader(w http.ResponseWriter, r *http.Request) error {
	result := helper.TVoidJsonResponse{Status: 200}
	//var posted TInputFile
	// json.NewDecoder(r.Body).Decode(&posted)

	file, fileHeader, _ := r.FormFile("file")
	fmt.Println(reflect.TypeOf(file), fileHeader.Filename)

	// fileName, err := helper.UploadS3(file, posted.FileName, []string{"test", "test"})
	// if err != nil {
	// 	fmt.Print(err)
	// }
	// fmt.Print(fileName)

	result.ResponseWrite(w)
	return nil
}
