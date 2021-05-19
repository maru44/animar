package helper

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func S3Session() *session.Session {
	creds := credentials.NewStaticCredentials(
		os.Getenv("S3_ACCESS_KEY_ID"),
		os.Getenv("S3_SECRET_KEY"),
		"",
	)
	sess := session.Must(session.NewSession(
		&aws.Config{
			Credentials: creds,
			Region:      aws.String("ap-northeast-1"),
		},
	))

	return sess
}

func UploadS3(file []byte, fileName string, pathList []string) (string, error) {
	contentType := getContentType(filepath.Ext(fileName))
	if contentType == "" {
		return "", errors.New("Unknown type")
	}

	sess := S3Session()
	u := s3manager.NewUploader(sess)
	slug := GenRandSlug(12)

	path := strings.Join(pathList, "/")
	key := fmt.Sprintf("%s/%s__%s", path, slug, fileName)

	_, err := u.Upload(&s3manager.UploadInput{
		Body:        bytes.NewReader(file),
		Bucket:      aws.String(os.Getenv("BUCKET")),
		ContentType: aws.String(contentType),
		Key:         aws.String(key),
	})
	if err != nil {
		return "", err
	}

	fileUrl := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", os.Getenv("BUCKET"), "ap-northeast-1", key)
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
