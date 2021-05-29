package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

func CloudStorageSession() *firebase.App {
	config := &firebase.Config{
		StorageBucket: os.Getenv("FIREBASE_STORAGE"),
	}
	opt := option.WithCredentialsFile("../../configs/secret_key.json")
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		fmt.Print(err)
	}
	return app
}

func CloudStorageClient() *storage.Client {
	app := CloudStorageSession()
	client, err := app.Storage(context.Background())
	if err != nil {
		fmt.Print(err)
	}
	return client
}

func UploadToCS(file, fileName string) {
	client := CloudStorageClient()
	bucket, err := client.DefaultBucket()
	if err != nil {
		fmt.Print(err)
	}

	contentType := getContentType(filepath.Ext(fileName))

	writer := bucket.Object(fileName).NewWriter(context.Background())
	writer.ObjectAttrs.ContentType = contentType
	writer.ObjectAttrs.CacheControl = "no-cache"

}
