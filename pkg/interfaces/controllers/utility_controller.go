package controllers

import (
	"animar/v1/pkg/tools/s3"
	"net/http"
)

type UtilityController struct {
	BaseController
}

func NewUtilityController() *UtilityController {
	return &UtilityController{
		BaseController: *NewBaseController(),
	}
}

func (controller *UtilityController) SimpleUploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB
	var returnFileName string

	file, fileHeader, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		returnFileName, err = s3.UploadS3(file, fileHeader.Filename, []string{"column", "content"})
	}
	response(w, err, map[string]interface{}{"data": returnFileName})
	return
}
