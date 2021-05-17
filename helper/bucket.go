package helper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"google.golang.org/grpc/internal/credentials"
)

func S3Session() {
	sess, err := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("", "", ""),
	}))
}
