package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
)

type GCSFileRepository struct {
	BucketName string
	Client     *storage.Client
}

func NewGCSFileRepository(client *storage.Client, bucketName string) *GCSFileRepository {
	return &GCSFileRepository{
		Client:     client,
		BucketName: bucketName,
	}
}

func (r *GCSFileRepository) UploadFile(file multipart.File, fileName string, contentType string) (string, error) {
	ctx := context.Background()
	writer := r.Client.Bucket(r.BucketName).Object(fileName).NewWriter(ctx)
	writer.ObjectAttrs.ContentType = contentType

	if _, err := io.Copy(writer, file); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", r.BucketName, fileName)
	return url, nil
}

func (r *GCSFileRepository) DeleteFile(fileURL string) error {
	ctx := context.Background()
	object := r.Client.Bucket(r.BucketName).Object(fileURL) // Use just object name in real use
	return object.Delete(ctx)
}
