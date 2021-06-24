package s3

import (
	"animar/v1/configs"
	"animar/v1/pkg/tools/tools"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func S3Session() *session.Session {
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

	return sess
}

func UploadS3(file multipart.File, fileName string, pathList []string) (string, error) {
	contentType := getContentType(filepath.Ext(fileName))
	if contentType == "" {
		return "", errors.New("Unknown type")
	}

	sess := S3Session()
	u := s3manager.NewUploader(sess)
	slug := tools.GenRandSlug(6)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		fmt.Print(err)
	}

	path := strings.Join(pathList, "/")
	key := fmt.Sprintf("%s/%s__%s", path, slug, fileName)

	_, err := u.Upload(&s3manager.UploadInput{
		Body:        buf,
		Bucket:      aws.String(configs.Bucket),
		ContentType: aws.String(contentType),
		Key:         aws.String(key),
	})
	if err != nil {
		return "", err
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
