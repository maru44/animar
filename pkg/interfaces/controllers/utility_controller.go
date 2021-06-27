package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/tools/s3"
	"net/http"
)

type UtilityController struct{}

func NewUtilityController() *UtilityController {
	return &UtilityController{}
}

func (controller *UtilityController) SimpleUploadImage(w http.ResponseWriter, r *http.Request) (ret error) {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB
	var returnFileName string
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		ret = response(w, domain.ErrUnauthorized, nil)
	} else {
		file, fileHeader, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			returnFileName, err = s3.UploadS3(file, fileHeader.Filename, []string{"column", "content"})
		}
		ret = response(w, err, map[string]interface{}{"data": returnFileName})
	}
	return ret
}
