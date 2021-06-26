package controllers

import (
	"animar/v1/pkg/domain"
	"animar/v1/pkg/infrastructure"
	"animar/v1/pkg/interfaces/apis"
	"animar/v1/pkg/tools/fire"
	"animar/v1/pkg/tools/s3"
	"net/http"
)

type UtilityController struct {
	api apis.ApiResponse
}

func NewUtilityController() *UtilityController {
	return &UtilityController{
		api: infrastructure.NewApiResponse(),
	}
}

func (controller *UtilityController) SimpleUploadImage(w http.ResponseWriter, r *http.Request) (ret error) {
	r.Body = http.MaxBytesReader(w, r.Body, 40*1024*1024) // 40MB
	var returnFileName string
	userId := fire.GetIdFromCookie(r)
	if userId == "" {
		ret = controller.api.Response(w, domain.ErrUnauthorized, nil)
	} else {
		file, fileHeader, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			returnFileName, err = s3.UploadS3(file, fileHeader.Filename, []string{"column", "content"})
		}
		ret = controller.api.Response(w, err, map[string]interface{}{"data": returnFileName})
	}
	return ret
}
