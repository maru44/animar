package infrastructure

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/interfaces/s3"
	"animar/v1/internal/pkg/tools/tools"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/maru44/perr"
)

type S3Uploader struct {
	Connect *s3manager.Uploader
}

func NewS3Uploader() s3.Uploader {
	creds := credentials.NewStaticCredentials(
		configs.S3AccessKeyId,
		configs.S3SecretKey,
		"",
	)
	sess := session.Must(session.NewSession(
		&aws.Config{
			Credentials: creds,
			Region:      aws.String("ap-northeast-1"),
		},
	))

	uploader := s3manager.NewUploader(sess)
	s3Uploader := new(S3Uploader)
	s3Uploader.Connect = uploader
	return s3Uploader
}

func (uploader *S3Uploader) ImageUploading(file multipart.File, fileName string, pathList []string) (string, error) {
	contentType := getContentType(filepath.Ext(fileName))
	if contentType == "" {
		return "", perr.New("Invalid file extension", perr.UnsupportedMediaType)
	}

	slug := tools.GenRandSlug(6)
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	path := strings.Join(pathList, "/")
	key := fmt.Sprintf("%s/%s__%s", path, slug, fileName)

	_, err := uploader.Connect.Upload(&s3manager.UploadInput{
		Body:        buf,
		Bucket:      aws.String(configs.Bucket),
		ContentType: aws.String(contentType),
		Key:         aws.String(key),
	})
	if err != nil {
		return "", perr.Wrap(err, perr.InternalServerErrorWithUrgency)
	}

	fileUrl := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", configs.Bucket, "ap-northeast-1", key)
	return fileUrl, nil
}

func getContentType(extension string) string {
	switch extension {
	case ".jpg":
		return "image/jpg"
	case ".jpeg":
		return "image/jpg"
	case ".png":
		return "image/png"
	default:
		return ""
	}
}
