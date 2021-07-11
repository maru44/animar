package domain

import "mime/multipart"

type S3Interactor interface {
	Image(multipart.File, string, []string) (string, error)
}
