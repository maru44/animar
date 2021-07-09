package s3

import "mime/multipart"

type Uploader interface {
	ImageUploading(multipart.File, string, []string) (string, error)
}
