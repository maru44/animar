package s3

import "mime/multipart"

type S3Repository struct {
	Uploader
}

func (repo *S3Repository) ImageUpload(file multipart.File, fileName string, pathList []string) (string, error) {
	returnFileName, err := repo.ImageUploading(file, fileName, pathList)
	return returnFileName, err
}
