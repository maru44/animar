package usecase

import (
	"animar/v1/internal/pkg/domain"
	"mime/multipart"
)

type S3Interactor struct {
	repository S3Repository
}

func NewS3Interactor(repo S3Repository) domain.S3Interactor {
	return &S3Interactor{
		repository: repo,
	}
}

/************************
        repository
************************/

type S3Repository interface {
	ImageUpload(multipart.File, string, []string) (string, error)
}

func (interactor *S3Interactor) Image(file multipart.File, fileName string, pathList []string) (string, error) {
	return interactor.repository.ImageUpload(file, fileName, pathList)
}
