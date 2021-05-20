package test

import (
	"animar/v1/helper"
	"fmt"
	"net/http"
)

type TInputFile struct {
	File     string `json:"file"`
	FileName string `json:"filename"`
}

func Uploader(w http.ResponseWriter, r *http.Request) error {
	result := helper.TVoidJsonResponse{Status: 200}

	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20MB
	file, fileHeader, _ := r.FormFile("file")
	file.Close()

	// この返値をdbに保存する
	_, err := helper.UploadS3(file, fileHeader.Filename, []string{"test", "2"})
	if err != nil {
		fmt.Print(err)
	}

	result.ResponseWrite(w)
	return nil
}
